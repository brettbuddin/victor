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
	Stop()
}

type Brain interface {
	Name() string
	Identity() User
	SetIdentity(User)
	Cacher
}

type Message interface {
	Body() string
	Room() Room
	User() User
	Reply(string) error
	SetParams([]string)
	Params() []string
}

type Room interface {
	Id() string
	Say(string) error
	Paste(string) error
	Sound(string) error
	Tweet(string) error
}

type User interface {
	Id() string
	Name() string
	AvatarURL() string
}

type CacheKeyer interface {
	CacheKey() string
}

type Cacher interface {
	Store(CacheKeyer)
	Get(string) CacheKeyer
	Exists(string) bool
	Expired(string) bool
	Delete(string)
}

func RoomKey(id string) string {
	return "room_" + id
}

func UserKey(id string) string {
	return "user_" + id
}
