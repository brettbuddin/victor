package campfire

import (
	"github.com/brettbuddin/campfire"
	"sync"
)

func NewCache() *cache {
	return &cache{
		mtx:   sync.RWMutex{},
		rooms: map[string]*campfire.Room{},
		users: map[string]*campfire.User{},
	}
}

type cache struct {
	mtx   sync.RWMutex
	rooms map[string]*campfire.Room
	users map[string]*campfire.User
}

func (c cache) Room(id string) *campfire.Room {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	if result, ok := c.rooms[id]; ok {
		return result
	}
	return nil
}

func (c cache) SetRoom(id string, room *campfire.Room) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.rooms[id] = room
}

func (c cache) User(id string) *campfire.User {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	if result, ok := c.users[id]; ok {
		return result
	}
	return nil
}

func (c cache) SetUser(id string, user *campfire.User) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.users[id] = user
}
