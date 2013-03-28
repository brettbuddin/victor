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
func (c *Campfire) Brain() *Brain {
    return c.brain
}

// Returns the Client used by this adapter.
func (c *Campfire) Client() *campfire.Client {
    return c.client
}

// Hear registers a new Hear Matcher with the Brain.
func (c *Campfire) Hear(expStr string, callback func(*Context)) {
    c.Brain().Hear(expStr, callback)
}

// Respond registers a new Respond Matcher with the Brain.
func (c *Campfire) Respond(expStr string, callback func(*Context)) {
    c.Brain().Respond(expStr, callback)
}

// Run is where the business happens.
func (c *Campfire) Run() {
    rooms := c.rooms
    messages := make(chan *campfire.Message)
    joined := 0

    for i := range rooms {
        me, err := c.Client().Me()

        if err != nil {
            log.Printf("Error fetching info about c: %s", err)
            continue
        }

        c.me = me

        room := c.Client().Room(rooms[i])

        if room.Join() != nil {
            log.Printf("Error joining room %i: %s", rooms[i], err)
            continue
        }
        joined++

        go c.pollRoomDetails(room)
        room.Stream(messages)
    }

    if joined == 0 {
        log.Fatal("No rooms joined; nothing to stream from.")
    }

    for msg := range messages {
        if msg.UserId() == c.me.Id() {
            continue
        }

        ctx := &Context{
            Reply: func(text string) {
                user := c.Brain().UserForId(msg.UserId())

                prefix := ""

                if user != nil {
                    prefix = user.Name() + ": "
                }

                c.Client().Room(msg.RoomId()).Say(prefix + text)
            },
            Send: func(text string) {
                c.Client().Room(msg.RoomId()).Say(text)
            },
            Paste: func(text string) {
                c.Client().Room(msg.RoomId()).Paste(text)
            },
            Sound: func(name string) {
                c.Client().Room(msg.RoomId()).Sound(name)
            },
        }

        go c.brain.Receive(ctx.SetMessage(msg))
    }
}

func (c *Campfire) pollRoomDetails(room *campfire.Room) {
    for {
        details, err := room.Show()

        if err != nil {
            continue
        }

        for _, user := range details.Users() {
            c.Brain().RememberUser(user)
        }

        time.Sleep(300 * time.Second)
    }
}
