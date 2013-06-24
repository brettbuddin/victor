package campfire

import (
    "github.com/brettbuddin/victor/adapter"
    "github.com/brettbuddin/campfire"
    "strconv"
    "fmt"
)

type Message struct {
    *campfire.Message
    user adapter.User
    room adapter.Room
    params []string
}

func (m *Message) Id() string {
    return strconv.Itoa(m.Message.Id)
}

func (m *Message) Body() string {
    return m.Message.Body
}

func (m *Message) User() adapter.User {
    return m.user
}

func (m *Message) Room() adapter.Room {
    return m.room
}

func (m *Message) Reply(text string) error {
    return m.Room().Say(fmt.Sprintf("%s: %s", m.User().Name(), text))
}

func (m *Message) SetParams(v []string) {
    m.params = v
}

func (m *Message) Params() []string {
    return m.params
}

type Room struct {
    *campfire.Room
}

func (r *Room) Id() string {
    return strconv.Itoa(r.Room.Id)
}

func (r *Room) Name() string {
    return r.Room.Name
}

func (r *Room) Say(text string) error {
    return r.Room.SendText(text)
}

func (r *Room) Paste(text string) error {
    return r.Room.SendPaste(text)
}

func (r *Room) Sound(name string) error {
    return r.Room.SendSound(name)
}

func (r *Room) Tweet(url string) error {
    return r.Room.SendTweet(url)
}

type User struct {
    *campfire.User
}

func (u *User) Id() string {
    return strconv.Itoa(u.User.Id)
}

func (u *User) Name() string {
    return u.User.Name
}
