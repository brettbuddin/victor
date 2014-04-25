package victor

import (
	"github.com/brettbuddin/victor/pkg/chat"
)

type Handler interface {
	Handle(*State)
}

type HandlerFunc func(*State)

func (f HandlerFunc) Handle(s *State) {
	f(s)
}

type State struct {
	robot   Robot
	message chat.Message
	params  []string
}

func (s *State) Robot() Robot {
	return s.robot
}

func (s *State) Chat() chat.Adapter {
	return s.robot.Chat()
}

func (s *State) Message() chat.Message {
	return s.message
}

func (s *State) Params() []string {
	return s.params
}
