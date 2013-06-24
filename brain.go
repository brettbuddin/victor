package victor

import (
    "github.com/brettbuddin/victor/adapter"
    "sync"
    "regexp"
    "strings"
    "log"
)

type Brain struct {
    mutex *sync.RWMutex
    name string
    id string
    users []adapter.User
    rooms []adapter.Room
    listeners []ListenerFunc
}

func NewBrain(name string) *Brain {
    return &Brain{
        mutex: &sync.RWMutex{},
        name: name,
        id: "",
        users: []adapter.User{},
        rooms: []adapter.Room{},
        listeners: []ListenerFunc{},
    }
}

func (b *Brain) Name() string {
    b.mutex.RLock()
    defer b.mutex.RUnlock()
    return b.name
}

func (b *Brain) Id() string {
    b.mutex.RLock()
    defer b.mutex.RUnlock()
    return b.id
}

func (b *Brain) SetId(v string) {
    b.mutex.Lock()
    defer b.mutex.Unlock()
    b.id = v
}

// Respond registers a listener that matches a direct message to
// the robot. This means "@botname command", "botname command" or
// "/command" forms
func (b *Brain) Respond(exp string, f func(adapter.Message)) (err error) {
    exp = strings.Join([]string{
        "(?i)",                                  // flags
        "\\A",                                   // begin
        "(?:(?:@)?" + b.Name() + "[:,]?\\s*|/)", // bot name
        "(?:" + exp + ")",                       // expression
        "\\z",                                   // end
    }, "")

    return b.Hear(exp, f)
}

// Hear registers a listener that matches any instance of the
// phrase in the room. Excluding from itself.
func (b *Brain) Hear(exp string, f func(adapter.Message)) (err error) {
    pattern, err := regexp.Compile(exp)

    if err != nil {
        return err
    }

    log.Printf("LISTENER: %s\n", exp)
    b.register(listenerFunc(pattern, f))
    return
}

func (b *Brain) register(l ListenerFunc) {
    b.mutex.Lock()
    defer b.mutex.Unlock()
    b.listeners = append(b.listeners, l)
}

// Receive accepts an incoming message and applies
// it to all listeners.
func (b *Brain) Receive(m adapter.Message) {
    b.mutex.RLock()
    defer b.mutex.RUnlock()
    for _, l := range b.listeners {
        go l(m)
    }
}

// AddUser caches a reference to the user for later lookup
func (b *Brain) AddUser(u adapter.User) {
    b.mutex.Lock()
    defer b.mutex.Unlock()
    b.users = append(b.users, u)
}

// AddRoom caches a reference to the room for later lookup
func (b *Brain) AddRoom(ro adapter.Room) {
    b.mutex.Lock()
    defer b.mutex.Unlock()
    b.rooms = append(b.rooms, ro)
}

// User returns the user with the specified ID
func (b *Brain) User(id string) adapter.User {
    b.mutex.RLock()
    defer b.mutex.RUnlock()
    for _, u := range b.users {
        if id == u.Id() {
            return u
        }
    }

    return nil
}

// Room returns the room with the specified ID
func (b *Brain) Room(id string) adapter.Room {
    b.mutex.RLock()
    defer b.mutex.RUnlock()
    for _, ro := range b.rooms {
        if id == ro.Id() {
            return ro
        }
    }

    return nil
}

type ListenerFunc func(adapter.Message)

func listenerFunc(pattern *regexp.Regexp, f func(adapter.Message)) ListenerFunc {
    return func(m adapter.Message) {
        results := pattern.FindAllStringSubmatch(m.Body(), -1)

        if len(results) > 0 {
            m.SetParams(results[0][1:])
            log.Printf("TRIGGER: %s\n", pattern)
            log.Printf("PARAMS: %s\n", m.Params())
            f(m)
        }
    }
}
