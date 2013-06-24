package adapter

import (
	"fmt"
)

var adapters = map[string]InitFunc{}

func Register(name string, init InitFunc) {
	adapters[name] = init
}

func Load(name string) (InitFunc, error) {
	a, ok := adapters[name]

	if !ok {
		return nil, fmt.Errorf("Unkown adapter %s", name)
	}

	return a, nil
}

type InitFunc func(Brain) Adapter

type Adapter interface {
	Listen(chan Message) error
}

type AdapterFunc func(chan Message)

func (f AdapterFunc) Produce(m chan Message) {
	f(m)
}

type Message interface {
	// Identity
	Id() string
	Body() string
	Room() Room
	User() User

	// Sends
	Reply(string) error

	// Params
	SetParams([]string)
	Params() []string
}

type Room interface {
	// Identity
	Id() string

	// Sends
	Say(string) error
	Paste(string) error
	Sound(string) error
	Tweet(string) error
}

type User interface {
	// Identity
	Id() string
	Name() string
}

type Brain interface {
	// Identity
	Id() string
	SetId(string)
	Name() string

	// Listener registration
	Respond(string, func(Message)) error
	Hear(string, func(Message)) error

	// Input
	Receive(Message)

	// Memory (Users and Rooms)
	AddUser(User)
	User(string) User
	UserExists(User) bool

	AddRoom(Room)
	Room(string) Room
	RoomExists(Room) bool
}
