package victor

import (
	"testing"
    "github.com/brettbuddin/victor/pkg/chat"
    "reflect"
)

var httpAddr = ":8000"

func init() {
    chat.Register("test", func(chat.Robot) chat.Adapter {
        return &testAdapter{}
    })
}

var routingTests = []struct{
    typ, pattern, text string
    matches []string
}{
    {"direct", "howdy", "ralph howdy", []string{}},
    {"direct", "tell (him|me)", "ralph tell him", []string{"him"}},
    {"direct", "build ([\\w-]+)/([\\w-]+):([\\w-]+)", "ralph build brettbuddin/victor:master", []string{"brettbuddin", "victor", "master"}},
    {"direct", "tell", "/tell", []string{}},
    {"indirect", "alot", "alot", []string{}},
}

func TestHandlers(t *testing.T) {
    robot := New(Config{
        Name: "ralph",
        ChatAdapter: "test",
        HTTPAddr: httpAddr,
    })

    go robot.Run()
    defer robot.Stop()

    for _, rt := range routingTests {
        if rt.typ == "direct" {
	        robot.HandleCommandFunc(rt.pattern, func(s State) {})
        } else {
	        robot.HandleFunc(rt.pattern, func(s State) {})
        }

	    handler, matches := robot.handler(&msg{text: rt.text})

        if handler == nil {
		    t.Errorf("No handler found for pattern=\"%s\": text=%s matches=%s", rt.pattern, rt.text, rt.matches)
            continue
        }

        if !reflect.DeepEqual(matches, rt.matches) {
		    t.Errorf("Incorrect matches for pattern=\"%s\": expected=%v got=%v", rt.pattern, rt.matches, matches)
        }
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

type testAdapter struct{}
func (a *testAdapter) Run() {}
func (a *testAdapter) Stop() {}
func (a *testAdapter) Send(channelID, msg string) {}
