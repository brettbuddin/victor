package brain

import (
	"github.com/brettbuddin/victor/pkg/adapter"
	"testing"
)

func TestRouting(t *testing.T) {
	b := New("ralph")
	called := 0

	b.Respond("howdy", func(m adapter.Message) {
		called++
	})
	b.Respond("tell (him|me)", func(m adapter.Message) {
		called++
	})
	b.Hear("alot", func(m adapter.Message) {
		called++
	})

	// Should trigger
	b.Receive(&msg{id: "123", body: "ralph howdy"})
	b.Receive(&msg{id: "123", body: "ralph tell him"})
	b.Receive(&msg{id: "123", body: "ralph tell me"})
	b.Receive(&msg{id: "123", body: "/tell me"})
	b.Receive(&msg{id: "123", body: "I heard alot of them."})

	if called != 5 {
		t.Error("One or more register actions weren't triggered")
	}
}

func TestParams(t *testing.T) {
	b := New("ralph")
	called := 0

	b.Respond("yodel (it)", func(m adapter.Message) {
		called++
		params := m.Params()
		if len(params) == 0 || params[0] != "it" {
			t.Errorf("Incorrect message params expected=%v got=%v", []string{"it"}, params)
		}
	})

	b.Receive(&msg{id: "123", body: "ralph yodel it"})

	if called != 1 {
		t.Error("Registered action was never triggered")
	}
}

func TestNonFiringRoutes(t *testing.T) {
	b := New("ralph")
	called := 0

	b.Respond("howdy", func(m adapter.Message) {
		called++
	})

	b.Receive(&msg{id: "123", body: "Tell ralph howdy."})

	if called > 0 {
		t.Error("Registered action was triggered when it shouldn't have been")
	}
}

type msg struct {
	id     string
	body   string
	user   adapter.User
	room   adapter.Room
	params []string
}

func (m *msg) Id() string {
	return m.id
}

func (m *msg) Body() string {
	return m.body
}

func (m *msg) Room() adapter.Room {
	return nil
}

func (m *msg) User() adapter.User {
	return nil
}

func (m *msg) Reply(s string) error {
	return nil
}

func (m *msg) SetParams(p []string) {
	m.params = p
}

func (m *msg) Params() []string {
	return m.params
}
