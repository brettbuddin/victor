package store

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
		return nil, fmt.Errorf("unkown adapter: %s", name)
	}

	return a, nil
}

type InitFunc func() Adapter

type Adapter interface {
	Get(string) (string, bool)
	Set(string, string)
	Delete(string)
	All() map[string]string
}
