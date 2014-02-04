package slack

import (
	"fmt"
	"github.com/brettbuddin/victor/pkg/chat"
	"github.com/gorilla/mux"
	"net/http"
	"regexp"
)

func init() {
	chat.Register("slack", func(r chat.Robot) chat.Adapter {
		return &slack{r}
	})
}

type slack struct {
	robot chat.Robot
}

func (s *slack) Run() {
	exp := regexp.MustCompile(fmt.Sprintf("\\A(?:@)?%s[:,]?\\s*", s.robot.Name))

	router := mux.NewRouter()
	router.HandleFunc("/slack/message", func(w http.ResponseWriter, r *http.Request) {
		text := r.PostFormValue("text")
		text = exp.ReplaceAllString(text, "")

		s.robot.Receive(&message{
			userId:      r.PostFormValue("user_id"),
			userName:    r.PostFormValue("user_name"),
			channelId:   r.PostFormValue("channel_id"),
			channelName: r.PostFormValue("channel_name"),
			text:        text,
		})
	}).Methods("POST")
}

func (s *slack) Send(channelId, msg string) {
}

func (s *slack) Stop() {
}

type message struct {
	userId      string `json:"user_id"`
	userName    string `json:"user_name"`
	channelId   string `json:"channel_id"`
	channelName string `json:"channel_name"`
	text        string `json:"text"`
}

func (m *message) UserId() string {
	return m.userId
}

func (m *message) UserName() string {
	return m.userName
}

func (m *message) ChannelId() string {
	return m.channelId
}

func (m *message) ChannelName() string {
	return m.channelName
}

func (m *message) Text() string {
	return m.text
}
