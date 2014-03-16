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

		roomIdStrings := strings.Split(roomsList, ",")
		rooms := []int{}

		for _, id := range roomIdStrings {
			j, err := strconv.Atoi(id)

			if err != nil {
				fmt.Printf("Room is not numeric: %s\n", id)
			}

			rooms = append(rooms, j)
		}

		return &adapter{
			robot:   r,
			account: account,
			token:   token,
			rooms:   rooms,
			stop:    make(chan bool),
			cache:   NewCache(),
		}
	})
}

type adapter struct {
	robot          chat.Robot
	account, token string
	rooms          []int
	stop           chan bool
	cache          *cache
}

func (a *adapter) Run() {
	client := campfire.NewClient(a.account, a.token)

	run := func(id int) {
		room, err := client.RoomForId(id)
		if err != nil {
			log.Printf("Unable to fetch room %d: %v\n", id, err)
			return
		}
		a.cache.SetRoom(strconv.Itoa(id), room)

		err = room.Join()
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
				userId := strconv.Itoa(msg.UserId)
				roomId := strconv.Itoa(msg.RoomId)
				room := a.cache.Room(roomId)

				var err error

				user := a.cache.User(userId)
				if user == nil {
					user, err = client.UserForId(msg.UserId)
					if err != nil {
						log.Printf("Unable to fetch user %d: %v\n", msg.UserId, err)
						continue
					}

					a.cache.SetUser(userId, user)
				}

				a.robot.Receive(&message{
					userId:   userId,
					userName: user.Name,
					roomId:   roomId,
					roomName: room.Name,
					text:     msg.Body,
				})
			}
		}
	}

	for _, id := range a.rooms {
		go run(id)
	}
}

func (a *adapter) Send(roomId, msg string) {
	room := a.cache.Room(roomId)
	if room == nil {
		log.Printf("Room %d hasn't been cached yet, unable to send message\n", roomId)
		return
	}

	err := room.SendText(msg)
	if err != nil {
		log.Printf("Error sending to room %d: %v\n", roomId, err)
	}
}

func (a *adapter) Stop() {
	close(a.stop)
	log.Println("Delaying shutdown by", SHUTDOWN_DELAY, "(for cleanup)")
	time.Sleep(SHUTDOWN_DELAY)
}

type message struct {
	userId, userName, roomId, roomName, text string
}

func (m *message) UserId() string {
	return m.userId
}

func (m *message) UserName() string {
	return m.userName
}

func (m *message) ChannelId() string {
	return m.roomId
}

func (m *message) ChannelName() string {
	return m.roomName
}

func (m *message) Text() string {
	return m.text
}
