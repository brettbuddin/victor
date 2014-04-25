package victor

import (
	"github.com/brettbuddin/victor/pkg/chat"
)

type Handler interface {
	Handle(State)
}

type HandlerFunc func(State)

func (f HandlerFunc) Handle(s State) {
	f(s)
}

type State interface {
	Robot() Robot
	Chat() chat.Adapter
	Message() chat.Message
	Params() []string
}

type state struct {
	robot   Robot
	message chat.Message
	params  []string
}

// Returns the Robot
func (s *state) Robot() Robot {
	return s.robot
}

// Returns the Chat adapter
func (s *state) Chat() chat.Adapter {
	return s.robot.Chat()
}

// Returns the Message
func (s *state) Message() chat.Message {
	return s.message
}

// Returns the params parsed from the Message
func (s *state) Params() []string {
	return s.params
}
