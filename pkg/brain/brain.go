package brain

import (
	"github.com/brettbuddin/victor/pkg/adapter"
	"github.com/brettbuddin/victor/pkg/cache"
	"log"
	"regexp"
	"strings"
	"time"
	"sync"
)

type ListenerFunc func(adapter.Message)

type Brain struct {
    *cache.Cache
	mutex     *sync.RWMutex
	name      string
	identity  adapter.User
	listeners []ListenerFunc
}

func New(name string) *Brain {
	return &Brain{
		Cache:     cache.New(3 * time.Hour),
		mutex:     &sync.RWMutex{},
		name:      name,
		listeners: []ListenerFunc{},
	}
}
func (b *Brain) Name() string {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.name
}

func (b *Brain) Identity() adapter.User {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.identity
}

func (b *Brain) SetIdentity(u adapter.User) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.identity = u
}

// Respond registers a listener that matches a direct message to
// the robot. This means "@botname command", "botname command" or
// "/command" forms
func (b *Brain) Respond(exp string, f func(adapter.Message)) (err error) {
	exp = strings.Join([]string{
		"(?i)", // flags
		"\\A",  // begin
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

	log.Printf("PATTERN: %s\n", exp)
	b.register(b.listener(pattern, f))
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
		l(m)
	}
}

func (b *Brain) listener(pattern *regexp.Regexp, f ListenerFunc) ListenerFunc {
	return func(m adapter.Message) {
		results := pattern.FindAllStringSubmatch(m.Body(), -1)

		if len(results) > 0 {
			m.SetParams(results[0][1:])
			log.Printf("MATCH=%s PARAMS=%s\n", pattern, m.Params())
			f(m)
		}
	}
}
