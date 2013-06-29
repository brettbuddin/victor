package shell

import (
	"github.com/brettbuddin/victor/adapter"
	"log"
	"strconv"
)

type Message struct {
	id     int
	body   string
	user   adapter.User
	room   adapter.Room
	params []string
}

func (m *Message) Id() string {
	return strconv.Itoa(m.id)
}

func (m *Message) Body() string {
	return m.body
}

func (m *Message) User() adapter.User {
	return m.user
}

func (m *Message) Room() adapter.Room {
	return m.room
}

func (m *Message) Reply(text string) error {
	log.Printf("REPLYING: %s: %s", m.User().Name(), text)
	return nil
}

func (m *Message) SetParams(v []string) {
	m.params = v
}

func (m *Message) Params() []string {
	return m.params
}

type Room struct {
	id   int
	name string
}

func (r Room) Id() string {
	return strconv.Itoa(r.id)
}

func (r Room) Name() string {
	return r.name
}

func (r Room) Say(text string) error {
	log.Println("SAYING:", text)
	return nil
}

func (r Room) Paste(text string) error {
	log.Println("PASTING:", text)
	return nil
}

func (r Room) Sound(name string) error {
	log.Println("PLAYING:", name)
	return nil
}

func (r Room) Tweet(url string) error {
	log.Println("DISPLAYING:", url)
	return nil
}

type User struct {
	id   int
	name string
}

func (u User) Id() string {
	return strconv.Itoa(u.id)
}

func (u User) Name() string {
	return u.name
}
