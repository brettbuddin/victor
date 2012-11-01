package victor

import (
    "github.com/brettbuddin/victor/campfire"
    "log"
    "time"
)

type Campfire struct {
    brain   *Brain
    account string
    token   string
    rooms   []int
    client  *campfire.Client
    me      *campfire.User
}

func NewCampfire(name string, account string, token string, rooms []int) *Campfire {
    return &Campfire{
        brain:   NewBrain(name),
        account: account,
        token:   token,
        rooms:   rooms,
        client:  campfire.NewClient(account, token),
    }
}

// Returns the Brain of this adapter.
func (self *Campfire) Brain() *Brain {
    return self.brain
}

// Returns the Client used by this adapter.
func (self *Campfire) Client() *campfire.Client {
    return self.client
}

// Hear registers a new Hear Matcher with the Brain.
func (self *Campfire) Hear(expStr string, callback func(*Context)) {
    self.Brain().Hear(expStr, callback)
}

// Respond registers a new Respond Matcher with the Brain.
func (self *Campfire) Respond(expStr string, callback func(*Context)) {
    self.Brain().Respond(expStr, callback)
}

// Run is where the business happens.
func (self *Campfire) Run() {
    rooms := self.rooms
    messages := make(chan *campfire.Message)
    joined := 0

    for i := range rooms {
        me, err := self.Client().Me()

        if err != nil {
            log.Printf("Error fetching info about self: %s", err)
            continue
        }

        self.me = me

        room := self.Client().Room(rooms[i])

        if room.Join() != nil {
            log.Printf("Error joining room %i: %s", rooms[i], err)
            continue
        }
        joined++

        go self.pollRoomDetails(room)
        room.Stream(messages)
    }

    if joined == 0 {
        log.Fatal("No rooms joined; nothing to stream from.")
    }

    for in := range messages {
        if in.UserId == self.me.Id {
            continue
        }

        if in.Type == "TextMessage" {
            ctx := &Context{
                Reply: func(text string) {
                    user := self.Brain().UserForId(in.UserId)

                    prefix := ""

                    if user != nil {
                        prefix = user.Name + ": "
                    }

                    self.Client().Room(in.RoomId).Say(prefix + text)
                },
                Send: func(text string) {
                    self.Client().Room(in.RoomId).Say(text)
                },
                Paste: func(text string) {
                    self.Client().Room(in.RoomId).Paste(text)
                },
                Sound: func(name string) {
                    self.Client().Room(in.RoomId).Sound(name)
                },
            }

            ctx.SetMessage(&Message{
                Id:        in.Id,
                Type:      in.Type,
                Body:      in.Body,
                CreatedAt: in.CreatedAt,
                UserId:    in.UserId,
                RoomId:    in.RoomId,
            })

            go self.brain.Receive(ctx)
        }
    }
}

func (self *Campfire) pollRoomDetails(room *campfire.Room) {
    for {
        details, err := room.Show()

        if err != nil {
            continue
        }

        for _, user := range details.Users {
            self.Brain().RememberUser(&User{Id: user.Id, Name: user.Name})
        }

        time.Sleep(300 * time.Second)
    }
}
