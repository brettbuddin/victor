package victor

import (
    "victor/campfire"
    "strings"
    "strconv"
    "log"
)

type Campfire struct {
    *Robot
    client *campfire.Client
}

func (self *Campfire) Run() {
    log.Print("Starting up...")

    rooms  := self.RoomList()
    client := self.Client() 

    channel := make(chan *campfire.Message)

    for i := range rooms {
        details, err := client.Room(rooms[i]).Show()

        if err != nil {
            continue
        }
        log.Print("Fetched room info.")

        for _, user := range details.Users {
            self.RememberUser(&User{Id: user.Id, Name: user.Name})
            log.Print("Remembering: " + user.Name)
        }

        room := client.Room(rooms[i])
        err   = room.Join()

        if err != nil {
            continue
        }
        log.Print("Joined room.")

        room.Stream(channel)
        log.Print("Listening...")
    }

    for {
        in   := <-channel

        if in.Type == "TextMessage" {
            msg := &TextMessage{
                Id: in.Id,
                Body: in.Body,
                CreatedAt: in.CreatedAt,

                Send: self.Send(in.RoomId),
                Reply: self.Reply(in.RoomId, in.UserId),
            }

            self.Receive(msg)
        }
    }
}

func (self *Campfire) Send(roomId int) func(string) {
    room := self.Client().Room(roomId)

    return func(text string) { 
        room.Say(text)
    }
}

func (self *Campfire) Reply(roomId int, userId int) func(string) {
    room   := self.Client().Room(roomId)
    user   := self.UserForId(userId)
    prefix := ""

    if user != nil {
       prefix = user.Name + ": "
    }

    return func(text string) { 
        room.Say(prefix + text)
    }
}

func (self *Campfire) Client() *campfire.Client {
    if self.client == nil {
        client := campfire.NewClient(
            self.Robot.options["account"], 
            self.Robot.options["token"],
        )
    
        self.client = client
    }
    
    return self.client
}

func (self *Campfire) RoomList() []int {
    if _, exists := self.Robot.options["rooms"]; !exists {
        panic("No rooms defined.")
    }

    roomIdStrings := strings.Split(self.Robot.options["rooms"], ",")
    rooms := make([]int, 0)

    for _, id := range roomIdStrings {
        j, _ := strconv.Atoi(id)
        rooms = append(rooms, j) 
    }

    return rooms
}
