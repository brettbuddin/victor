package victor

import (
    "github.com/brettbuddin/victor/campfire"
    "strings"
    "strconv"
    "log"
    "os"
)

type Campfire struct {
    *Robot
    account string
    token   string
    rooms   []int
    client *campfire.Client
}

func NewCampfire(robot *Robot) *Campfire {
    c := &Campfire{Robot: robot}
    
    c = loadEnv(c)

    return c
}

func loadEnv(c *Campfire) *Campfire {
    account := os.Getenv("VICTOR_CAMPFIRE_ACCOUNT")
    token   := os.Getenv("VICTOR_CAMPFIRE_TOKEN")
    rooms   := os.Getenv("VICTOR_CAMPFIRE_ROOMS")

    if account == "" {
        log.Panic("No account set.")
    }

    if token == "" {
        log.Panic("No token set.")
    }

    if rooms == "" {
        log.Panic("No rooms set.")
    }

    c.account = account
    c.token   = token
    c.client  = campfire.NewClient(account, token)

    roomIdStrings := strings.Split(rooms, ",")
    roomsArr      := make([]int, 0)

    for _, id := range roomIdStrings {
        j, _ := strconv.Atoi(id)
        roomsArr = append(roomsArr, j) 
    }

    c.rooms  = roomsArr

    return c
}


func (self *Campfire) Run() {
    log.Print("Starting up...")

    rooms  := self.rooms

    channel := make(chan *campfire.Message)

    for i := range rooms {
        details, err := self.client.Room(rooms[i]).Show()

        if err != nil {
            log.Printf("Error fetching room info %i: %s", rooms[i], err)
            continue
        }
        log.Print("Fetched room info.")

        for _, user := range details.Users {
            self.RememberUser(&User{Id: user.Id, Name: user.Name})
            log.Print("Remembering: " + user.Name)
        }

        room := self.client.Room(rooms[i])
        err   = room.Join()

        if err != nil {
            log.Printf("Error joining room %i: %s", rooms[i], err)
            continue
        }
        log.Print("Joined room.")

        room.Stream(channel)
        log.Print("Listening...")
    }

    for {
        in := <-channel

        if in.Type == "TextMessage" {
            msg := &TextMessage{
                Id: in.Id,
                Body: in.Body,
                CreatedAt: in.CreatedAt,

                Send: self.Send(in.RoomId),
                Reply: self.Reply(in.RoomId, in.UserId),
            }

            go self.Receive(msg)
        }
    }
}

func (self *Campfire) Send(roomId int) func(string) {
    room := self.client.Room(roomId)

    return func(text string) { 
        room.Say(text)
    }
}

func (self *Campfire) Reply(roomId int, userId int) func(string) {
    room   := self.client.Room(roomId)
    user   := self.UserForId(userId)
    prefix := ""

    if user != nil {
       prefix = user.Name + ": "
    }

    return func(text string) { 
        room.Say(prefix + text)
    }
}
