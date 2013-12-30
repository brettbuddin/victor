package victor

import (
	"github.com/brettbuddin/victor/pkg/adapter"
	"github.com/brettbuddin/victor/pkg/brain"
	_ "github.com/brettbuddin/victor/pkg/adapter/campfire"
	_ "github.com/brettbuddin/victor/pkg/adapter/shell"
	"log"
)

type Robot struct {
	adapter  adapter.Adapter
	brain    *brain.Brain
	incoming chan adapter.Message
	stop     chan bool
}

type Message interface {
    adapter.Message
}

// New returns a Robot
func New(adapterName, robotName string) (*Robot, error) {
	initFunc, err := adapter.Load(adapterName)

	if err != nil {
		return nil, err
	}

	brain := brain.New(robotName)
	bot := &Robot{
		adapter:  initFunc(brain),
		brain:    brain,
		incoming: make(chan adapter.Message),
		stop:     make(chan bool),
	}

	defaults(bot)
	return bot, nil
}

func (r *Robot) Brain() *brain.Brain {
	return r.brain
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
	go func() {
	    err := r.adapter.Listen(r.incoming)

	    if err != nil {
            log.Println(err)
	    }

        r.Stop()
	}()

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
	close(r.stop)
}
