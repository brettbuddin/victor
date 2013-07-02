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
type AdapterFunc func(chan Message) error

func (f AdapterFunc) Listen(m chan Message) error {
	return f(m)
}

type Adapter interface {
	Listen(chan Message) error
}

type Brain interface {
	Name() string
	Identity() User
	SetIdentity(User)
	Cache() Cacher
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

type Cacheable interface {
	CacheKey() string
}

type Cacher interface {
	Add(Cacheable)
	Get(string) Cacheable
	Exists(string) bool
	Delete(string)
}

func RoomKey(id string) string {
	return "room_" + id
}

func UserKey(id string) string {
	return "user_" + id
}
