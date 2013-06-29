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

type InitFunc func(Agent) Adapter
type AdapterFunc func(chan Message)

func (f AdapterFunc) Produce(m chan Message) {
	f(m)
}

type Adapter interface {
	Listen(chan Message) error
}

type Agent interface {
	Identity() User
	SetIdentity(User)
}

type Message interface {
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
	Id() string

	// Sends
	Say(string) error
	Paste(string) error
	Sound(string) error
	Tweet(string) error
}

type User interface {
	Id() string
	Name() string
}
