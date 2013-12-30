package campfire

import (
	"errors"
	"github.com/brettbuddin/campfire"
	"github.com/brettbuddin/victor/pkg/adapter"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	SHUTDOWN_DELAY = 5 * time.Second
)

func init() {
	adapter.Register("campfire", func(b adapter.Brain) adapter.Adapter {
		account := os.Getenv("VICTOR_CAMPFIRE_ACCOUNT")
		token := os.Getenv("VICTOR_CAMPFIRE_TOKEN")
		roomList := os.Getenv("VICTOR_CAMPFIRE_ROOMS")

		if account == "" || token == "" || roomList == "" {
			log.Println("The following environment variables are required:")
			log.Println("VICTOR_CAMPFIRE_ACCOUNT, VICTOR_CAMPFIRE_TOKEN, VICTOR_CAMPFIRE_ROOMS")
			os.Exit(1)
		}

		client := campfire.NewClient(account, token)
		roomIdStrings := strings.Split(roomList, ",")
		roomIds := []int{}

		for _, id := range roomIdStrings {
			j, err := strconv.Atoi(id)

			if err != nil {
				log.Printf("Room is not numeric: %s\n", id)
			}

			roomIds = append(roomIds, j)
		}

		return &Adapter{
			brain:   b,
			client:  client,
			stop:    make(chan bool),
			roomIds: roomIds,
		}
	})
}

type Adapter struct {
	brain   adapter.Brain
	client  *campfire.Client
	stop    chan bool
	roomIds []int
}

func (a *Adapter) Listen(outgoing chan adapter.Message) error {
	rooms, err := a.bootstrap()
	if err != nil {
		return err
	}

	streams := []*campfire.Stream{}
	incoming := make(chan *campfire.Message)

	for _, room := range rooms {
		s := room.Stream()
		go s.Connect()
		go func() {
			for {
				message, ok := <-s.Messages()
				if !ok {
					break
				}

				incoming <- message
			}
		}()
		streams = append(streams, s)
	}

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

			close(incoming)
			return nil
		case m := <-incoming:
			roomId := itoa(m.RoomId)
			userId := itoa(m.UserId)

			key := adapter.UserKey(userId)

			if !a.brain.Exists(key) || a.brain.Expired(key) {
				user, err := a.client.UserForId(m.UserId)

				if err != nil {
					break
				}

				a.brain.Store(User{user})
			}

			outgoing <- &Message{
				message: m,
				room:    a.brain.Get(adapter.RoomKey(roomId)).(Room),
				user:    a.brain.Get(adapter.UserKey(userId)).(User),
			}
		}
	}
}

func (a *Adapter) Stop() {
	close(a.stop)
	log.Println("Delaying shutdown by", SHUTDOWN_DELAY, "(for cleanup)")
	time.Sleep(SHUTDOWN_DELAY)
}

func (a *Adapter) bootstrap() ([]*campfire.Room, error) {
	rooms := []*campfire.Room{}

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

		a.brain.Store(Room{room})
		for _, u := range room.Users {
			a.brain.Store(User{u})
		}

		log.Printf("ROOM[%d]: joined %s\n", id, room.Name)
		rooms = append(rooms, room)
	}

	if len(rooms) == 0 {
		return nil, errors.New("No rooms joined")
	}

	return rooms, nil
}
