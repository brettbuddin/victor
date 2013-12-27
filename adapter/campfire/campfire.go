package campfire

import (
	"errors"
	"github.com/brettbuddin/campfire"
	"github.com/brettbuddin/victor/adapter"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

func init() {
	adapter.Register("campfire", func(b adapter.Brain) adapter.Adapter {
		return &Adapter{
			brain: b,
			once: sync.Once{},
			stop: make(chan bool),
			roomIds: []int{},
		}
	})
}

type Adapter struct {
	brain  adapter.Brain
	client *campfire.Client
	once   sync.Once
	stop   chan bool
	roomIds []int
}

func (a *Adapter) Listen(messages chan adapter.Message) error {
	a.configure()

	rooms, err := a.bootstrap()
	if err != nil {
		return err
	}

	streams := []*campfire.Stream{}
	rawMessages := make(chan *campfire.Message)

	for _, room := range rooms {
        s := room.Stream(rawMessages)
        go s.Connect()
        streams = append(streams, s)
	}

	cache := a.brain.Cache()

    for {
        select {
        case <-a.stop:
            log.Println("Disconnecting from streams")
            for _, s := range streams {
                s.Disconnect()
            }

            log.Println("Leaving rooms")
            for _, r := range rooms {
                r.Leave()
            }

            close(messages)
            return nil
        case m := <-rawMessages:
            roomId := itoa(m.RoomId)
            userId := itoa(m.UserId)

            if !cache.Exists(adapter.UserKey(userId)) {
                user, err := a.client.UserForId(m.UserId)

                if err != nil {
                    break
                }

                cache.Add(User{user})
            }

            messages <- &Message{
                message: m,
                room:    cache.Get(adapter.RoomKey(roomId)).(Room),
                user:    cache.Get(adapter.UserKey(userId)).(User),
            }
        }
    }
}

func (a *Adapter) Stop() {
    a.stop <- true
    close(a.stop)
}

func (a *Adapter) configure() {
    a.once.Do(func() {
        account := os.Getenv("VICTOR_CAMPFIRE_ACCOUNT")
        token := os.Getenv("VICTOR_CAMPFIRE_TOKEN")
        roomList := os.Getenv("VICTOR_CAMPFIRE_ROOMS")

        if account == "" || token == "" || roomList == "" {
            log.Println("The following environment variables are required:")
            log.Println("VICTOR_CAMPFIRE_ACCOUNT, VICTOR_CAMPFIRE_TOKEN, VICTOR_CAMPFIRE_ROOMS")
            os.Exit(1)
        }

        a.client = campfire.NewClient(account, token)
        roomIdStrings := strings.Split(roomList, ",")

        for _, id := range roomIdStrings {
            j, err := strconv.Atoi(id)

            if err != nil {
                log.Printf("Room is not numeric: %s\n", id)
            }

            a.roomIds = append(a.roomIds, j)
        }
	})
}

func (a *Adapter) bootstrap() ([]*campfire.Room, error) {
	rooms := []*campfire.Room{}
	cache := a.brain.Cache()

	me, err := a.client.Me()
	if err != nil {
		return nil, errors.New("Could not get information about self user")
	}

	a.brain.SetIdentity(User{me})

	for _, id := range a.roomIds {
		room, err := a.client.RoomForId(id)

		if err != nil {
			log.Printf("ROOM[%d]: unable to get info on room: %s\n", id, err)
			continue
		}

		err = room.Join()

		if err != nil {
			log.Printf("ROOM[%d]: Unable to join room\n", id)
			continue
		}

		cache.Add(Room{room})
		for _, u := range room.Users {
			cache.Add(User{u})
		}

		log.Printf("ROOM[%d]: joined %s\n", id, room.Name)
		rooms = append(rooms, room)
	}

	if len(rooms) == 0 {
		return nil, errors.New("No rooms joined")
	}

	return rooms, nil
}
