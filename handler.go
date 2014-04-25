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

func (s *state) Robot() Robot {
	return s.robot
}

func (s *state) Chat() chat.Adapter {
	return s.robot.Chat()
}

func (s *state) Message() chat.Message {
	return s.message
}

func (s *state) Params() []string {
	return s.params
}
