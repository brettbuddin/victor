package campfire

import (
	"github.com/brettbuddin/campfire"
	"github.com/brettbuddin/victor/adapter"
	"log"
	"os"
	"strconv"
	"strings"
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
		roomIds := make([]int, 0)

		for _, id := range roomIdStrings {
			j, err := strconv.Atoi(id)

			if err != nil {
				log.Fatalf("room is not numeric: %s\n", id)
			}

			roomIds = append(roomIds, j)
		}

		return &Campfire{
			Brain:   b,
			client:  client,
			roomIds: roomIds,
		}
	})
}

type Campfire struct {
	adapter.Brain
	client  *campfire.Client
	roomIds []int
}

func (c *Campfire) Listen(messages chan adapter.Message) (err error) {
	me, err := c.client.Me()

	if err != nil {
		log.Fatalf("CAMPFIRE: could not fetch info about self: %s", err)
	}

	c.Brain.SetId(strconv.Itoa(me.Id))

	rooms := []*campfire.Room{}

	for _, id := range c.roomIds {
		room, err := c.client.RoomForId(id)

		if err != nil {
			log.Printf("ROOM[%d]: unable to get info on room: %s\n", id, err)
			continue
		}

		err = room.Join()

		if err != nil {
			log.Printf("ROOM[%d]: Unable to join room\n", id)
			continue
		}

		c.Brain.AddRoom(&Room{room})
		for _, u := range room.Users {
			c.Brain.AddUser(&User{u})
		}

		log.Printf("ROOM[%d]: joined %s\n", id, room.Name)
		rooms = append(rooms, room)
	}

	if len(rooms) == 0 {
		log.Fatal("CAMPFIRE: No rooms joined")
	}

	rawMessages := make(chan *campfire.Message)

	for _, room := range rooms {
		go room.Stream(rawMessages).Connect()
	}

	for {
		select {
		case m := <-rawMessages:
			messages <- &Message{
				Message: m,
				room:    c.Brain.Room(strconv.Itoa(m.RoomId)),
				user:    c.Brain.User(strconv.Itoa(m.UserId)),
			}
		}
	}
}
