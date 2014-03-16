package campfire

import (
	"fmt"
	"github.com/brettbuddin/campfire"
	"github.com/brettbuddin/victor/pkg/chat"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const SHUTDOWN_DELAY = 5 * time.Second

func init() {
	chat.Register("campfire", func(r chat.Robot) chat.Adapter {
		roomsList := os.Getenv("VICTOR_CAMPFIRE_ROOMS")
		account := os.Getenv("VICTOR_CAMPFIRE_ACCOUNT")
		token := os.Getenv("VICTOR_CAMPFIRE_TOKEN")

		if roomsList == "" || account == "" || token == "" {
			log.Println("The following environment variables are required:")
			log.Println("VICTOR_CAMPFIRE_ACCOUNT, VICTOR_CAMPFIRE_TOKEN, VICTOR_CAMPFIRE_ROOMS")
			os.Exit(1)
		}

		roomIDStrings := strings.Split(roomsList, ",")
		roomIDs := []int{}

		for _, id := range roomIDStrings {
			j, err := strconv.Atoi(id)

			if err != nil {
				fmt.Printf("Room is not numeric: %s\n", id)
			}

			roomIDs = append(roomIDs, j)
		}

		return &adapter{
			robot:   r,
			client:  campfire.NewClient(account, token),
			roomIDs: roomIDs,
			stop:    make(chan struct{}),
		}
	})
}

type adapter struct {
	robot   chat.Robot
	client  *campfire.Client
	roomIDs []int
	stop    chan struct{}
}

func (a *adapter) Run() {
	run := func(id int) {
		room := campfire.Room{
			Connection: a.client.Connection,
			ID:         id,
		}

		err := room.Join()
		if err != nil {
			log.Printf("Unable to join room %d: %v\n", id, err)
			return
		}

		stream := room.Stream()
		go stream.Connect()
		messages := stream.Messages()

		for {
			select {
			case <-a.stop:
				stream.Disconnect()
				room.Leave()
				return
			case msg := <-messages:
				userID := strconv.Itoa(msg.UserID)
				roomID := strconv.Itoa(msg.RoomID)

				roomName, ok := a.getStoreKey("room.name." + roomID)
				if !ok {
					room, err := a.client.RoomForID(id)
					if err != nil {
						log.Printf("Unable to fetch room %d: %v\n", id, err)
						return
					}

					roomName = room.Name
					a.setStoreKey("room.name."+roomID, room.Name)
				}

				userName, ok := a.getStoreKey("user.name." + userID)
				if !ok {
					user, err := a.client.UserForID(msg.UserID)
					if err != nil {
						log.Printf("Unable to fetch user %d: %v\n", msg.UserID, err)
						continue
					}

					userName = user.Name
					a.setStoreKey("user.name."+userID, user.Name)
				}

				a.robot.Receive(&message{
					userID:   userID,
					userName: userName,
					roomID:   roomID,
					roomName: roomName,
					text:     msg.Body,
				})
			}
		}
	}

	for _, id := range a.roomIDs {
		go run(id)
	}
}

func (a *adapter) Send(roomID, msg string) {
	id, _ := strconv.Atoi(roomID)
	room := campfire.Room{
		Connection: a.client.Connection,
		ID:         id,
	}

	err := room.SendText(msg)
	if err != nil {
		log.Printf("Error sending to room %d: %v\n", roomID, err)
	}
}

func (a *adapter) Stop() {
	close(a.stop)
	log.Println("Delaying shutdown by", SHUTDOWN_DELAY, "(for cleanup)")
	time.Sleep(SHUTDOWN_DELAY)
}

func (a *adapter) getStoreKey(key string) (string, bool) {
	return a.robot.Store().Get("campfire." + key)
}

func (a *adapter) setStoreKey(key, val string) {
	a.robot.Store().Set("campfire."+key, val)
}

type message struct {
	userID, userName, roomID, roomName, text string
}

func (m *message) UserID() string {
	return m.userID
}

func (m *message) UserName() string {
	return m.userName
}

func (m *message) ChannelID() string {
	return m.roomID
}

func (m *message) ChannelName() string {
	return m.roomName
}

func (m *message) Text() string {
	return m.text
}

func key(key string) string {
	return "campfire." + key
}
