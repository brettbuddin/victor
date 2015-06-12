package store

import (
	"fmt"

	"github.com/gorilla/mux"
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

type Robot interface {
	Name() string
	HTTP() *mux.Router
	StoreConfig() (interface{}, bool)
}

type Adapter interface {
	Get(string) (string, bool)
	Set(string, string)
	Delete(string)
	All() map[string]string
}
