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
    Id()   string
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
}

type Brain interface {
    Name() string
    Id() string
    SetId(string)
    Receive(Message)
    Respond(string, func(Message)) error
    Hear(string, func(Message)) error

    AddUser(User)
    AddRoom(Room)
    User(string) User
    Room(string) Room
}
