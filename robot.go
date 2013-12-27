package victor

import (
	"github.com/brettbuddin/victor/adapter"
	_ "github.com/brettbuddin/victor/adapter/campfire"
	_ "github.com/brettbuddin/victor/adapter/shell"
	"log"
)

type Robot struct {
	adapter  adapter.Adapter
	brain    *Brain
	incoming chan adapter.Message
	stop     chan bool
}

type Message interface {
	adapter.Message
}

type ImmutableBrain interface {
	Name() string
	Identity() adapter.User
	Cache() adapter.Cacher
}

// New returns a Robot
func New(adapterName, robotName string) (*Robot, error) {
	initFunc, err := adapter.Load(adapterName)

	if err != nil {
		return nil, err
	}

	brain := NewBrain(robotName)
	bot := &Robot{
		adapter:  initFunc(brain),
		brain:    brain,
		stop:     make(chan bool),
		incoming: make(chan adapter.Message),
	}

	defaults(bot)
	return bot, nil
}

func (r *Robot) Brain() ImmutableBrain {
	return ImmutableBrain(r.brain)
}

// Respond proxies the registration of a respond
// command to the brain.
func (r *Robot) Respond(exp string, f func(Message)) (err error) {
	return r.brain.Respond(exp, func(m adapter.Message) {
		f(m.(Message))
	})
}

// Hear proxies the registration of a hear
// command to the brain.
func (r *Robot) Hear(exp string, f func(Message)) (err error) {
	return r.brain.Hear(exp, func(m adapter.Message) {
		f(m.(Message))
	})
}

// Run starts the robot.
func (r *Robot) Run() error {
	go r.adapter.Listen(r.incoming)

	for {
		select {
		case <-r.stop:
			log.Println("Stopping")
			r.adapter.Stop()
			return nil
		case m := <-r.incoming:
			if r.brain.Identity() == nil || m.User().Id() != r.brain.Identity().Id() {
				go r.brain.Receive(m)
			}
		}
	}
}

func (r *Robot) Stop() {
	r.stop <- true
	close(r.stop)
}
