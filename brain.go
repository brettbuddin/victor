package victor

import (
    "log"
    "sync"
    "regexp"
)

type Adapter interface {
    Brain() *Brain
    Hear(string, func(*Context))
    Respond(string, func(*Context))
    Run()
}

type Brain struct {
    name     string
    options  map[string]string
    matchers []*Matcher

    userLock *sync.Mutex
    users    []User
}

func NewBrain(name string) *Brain {
    brain := &Brain{
        name:     name,
        matchers: make([]*Matcher, 0, 1),
        userLock: &sync.Mutex{},
    }

    registerDefaultAbilities(brain)

    return brain
}

// AddMatcher adds a Matcher to the Brain's list of matching patterns.
func (b *Brain) AddMatcher(m *Matcher) {
    b.matchers = append(b.matchers, m)

    log.Printf("Pattern: /%s/", m.Pattern.String())
}

// Matchers returns the list of known matching patterns.
func (b *Brain) Matchers() []*Matcher {
    return b.matchers
}

// Hear creates and registers a new Matcher with the Brain that is triggered
// when the pattern matches anything said in the room.
func (b *Brain) Hear(expStr string, callback func(*Context)) {
    exp, _ := regexp.Compile(expStr)

    b.AddMatcher(NewMatcher(exp, callback))
}

// Respond creates and registers a new Matcher with the Brain that is triggered
// when the pattern matches a statement directed at the bot specifically.
func (b *Brain) Respond(expStr string, callback func(*Context)) {
    expWithNameStr := "(?i)\\A(?:" + b.name + "[:,]?\\s*|/)(?:" + expStr + ")\\z"
    exp, _ := regexp.Compile(expWithNameStr)

    b.AddMatcher(NewMatcher(exp, callback))
}

// Receive takes input from the service adapter and tests it against
// all registered Matchers. Breaks on the first match.
func (b *Brain) Receive(ctx *Context) {
    for _, matcher := range b.matchers {
        if matcher.Test(ctx) {
            matcher.Callback(ctx)
            return
        }
    }
}

// Instructs the Brain to remember the information about a user.
// Passing it a user that it has already seen will update the info
// for the same user in memory.
func (b *Brain) RememberUser(user User) {
    b.userLock.Lock()
    defer b.userLock.Unlock()

    for i, u := range b.users {
        if u.Id() == user.Id() {
            // update the name if its different
            if u.Name() != user.Name() {
                b.users[i] = user
            }

            return
        }
    }

    b.users = append(b.users, user)
}

// Returns a list of all known users.
func (b *Brain) Users() []User {
    return b.users
}

// Returns the User with the specified ID from memory.
func (b *Brain) UserForId(id int) User {
    b.userLock.Lock()
    defer b.userLock.Unlock()

    for _, user := range b.users {
        if user.Id() == id {
            return user
        }
    }

    return nil
}
