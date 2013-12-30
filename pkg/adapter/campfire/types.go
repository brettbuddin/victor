package campfire

import (
	"fmt"
	"github.com/brettbuddin/campfire"
	"github.com/brettbuddin/victor/pkg/adapter"
	"strconv"
)

type Message struct {
	message *campfire.Message
	user    adapter.User
	room    adapter.Room
	params  []string
}

func (m *Message) Id() string {
	return itoa(m.message.Id)
}

func (m *Message) Body() string {
	return m.message.Body
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
	room *campfire.Room
}

func (r Room) Id() string {
	return itoa(r.room.Id)
}

func (r Room) CacheKey() string {
	return adapter.RoomKey(r.Id())
}

func (r Room) Name() string {
	return r.room.Name
}

func (r Room) Say(text string) error {
	return r.room.SendText(text)
}

func (r Room) Paste(text string) error {
	return r.room.SendPaste(text)
}

func (r Room) Sound(name string) error {
	return r.room.SendSound(name)
}

func (r Room) Tweet(url string) error {
	return r.room.SendTweet(url)
}

type User struct {
	user *campfire.User
}

func (u User) CacheKey() string {
	return adapter.UserKey(u.Id())
}

func (u User) Id() string {
	return itoa(u.user.Id)
}

func (u User) Name() string {
	return u.user.Name
}

func (u User) AvatarURL() string {
	return u.user.AvatarURL
}

func itoa(i int) string {
	return strconv.Itoa(i)
}
