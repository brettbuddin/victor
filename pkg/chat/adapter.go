package chat

import (
	"fmt"
	"github.com/brettbuddin/victor/pkg/httpserver"
	"github.com/brettbuddin/victor/pkg/store"
)

var adapters = map[string]InitFunc{}

func Register(name string, init InitFunc) {
	adapters[name] = init
}

func Load(name string) (InitFunc, error) {
	a, ok := adapters[name]

	if !ok {
		return nil, fmt.Errorf("unkown adapter: %s", name)
	}

	return a, nil
}

type InitFunc func(Robot) Adapter

type Adapter interface {
	Run()
	Send(string, string)
	Stop()
}

type Robot interface {
	Name() string
	HTTP() *httpserver.Server
	Store() store.Store
	Chat() Adapter
	Receive(Message)
}

type Message interface {
	UserId() string
	UserName() string
	ChannelId() string
	ChannelName() string
	Text() string
}
