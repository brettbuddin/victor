package victor

import (
	"github.com/brettbuddin/victor/adapter"
	_ "github.com/brettbuddin/victor/adapter/campfire"
	_ "github.com/brettbuddin/victor/adapter/shell"
)

type Robot struct {
	adapter adapter.Adapter
	brain   adapter.Brain
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

	brain := NewBrain(robotName)
	bot := &Robot{
		adapter: initFunc(brain),
		brain:   brain,
	}

	defaults(bot)
	return bot, nil
}

// Brain returns the brain of the robot
func (r *Robot) Brain() adapter.Brain {
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
	messages := make(chan adapter.Message)
	done := make(chan bool)

	go func() {
		r.adapter.Listen(messages)
		done <- true
	}()

	for {
		select {
		case <-done:
			close(done)
			close(messages)
			return nil
		case m := <-messages:
			if m.User().Id() != r.brain.Id() {
				r.brain.Receive(m)
			}
		}
	}

	return nil
}
