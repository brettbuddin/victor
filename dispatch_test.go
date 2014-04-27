package victor

import (
	"testing"
)

func TestRouting(t *testing.T) {
	robot := &robot{name: "ralph"}
	robot.dispatch = newDispatch(robot)

	called := 0

	robot.HandleCommandFunc("howdy", func(s State) {
		called++
	})
	robot.HandleCommandFunc("tell (him|me)", func(s State) {
		called++
	})
	robot.HandleFunc("alot", func(s State) {
		called++
	})

	// Should trigger
	robot.ProcessMessage(&msg{text: "ralph howdy"})
	robot.ProcessMessage(&msg{text: "ralph tell him"})
	robot.ProcessMessage(&msg{text: "ralph tell me"})
	robot.ProcessMessage(&msg{text: "/tell me"})
	robot.ProcessMessage(&msg{text: "I heard alot of them."})

	if called != 5 {
		t.Errorf("One or more register actions weren't triggered")
	}
}

func TestParams(t *testing.T) {
	robot := &robot{name: "ralph"}
	robot.dispatch = newDispatch(robot)

	called := 0

	robot.HandleCommandFunc("yodel (it)", func(s State) {
		called++
		params := s.Params()
		if len(params) == 0 || params[0] != "it" {
			t.Errorf("Incorrect message params expected=%v got=%v", []string{"it"}, params)
		}
	})

	robot.ProcessMessage(&msg{text: "ralph yodel it"})

	if called != 1 {
		t.Error("Registered action was never triggered")
	}
}

func TestNonFiringRoutes(t *testing.T) {
	robot := &robot{name: "ralph"}
	robot.dispatch = newDispatch(robot)

	called := 0

	robot.HandleCommandFunc("howdy", func(s State) {
		called++
	})

	robot.ProcessMessage(&msg{text: "Tell ralph howdy."})

	if called > 0 {
		t.Error("Registered action was triggered when it shouldn't have been")
	}
}

type msg struct {
	userID      string
	userName    string
	channelID   string
	channelName string
	text        string
}

func (m *msg) UserID() string {
	return m.userID
}

func (m *msg) UserName() string {
	return m.userName
}

func (m *msg) ChannelID() string {
	return m.channelID
}

func (m *msg) ChannelName() string {
	return m.channelName
}

func (m *msg) Text() string {
	return m.text
}
