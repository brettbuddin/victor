package slack

import (
	"encoding/json"
	"fmt"
	"github.com/brettbuddin/victor/pkg/chat"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/url"
	"os"
)

func init() {
	chat.Register("slack", func(r chat.Robot) chat.Adapter {
		return &slack{
			robot: r,
			team:  os.Getenv("VICTOR_SLACK_TEAM"),
			token: os.Getenv("VICTOR_SLACK_TOKEN"),
		}
	})
}

type slack struct {
	robot chat.Robot
	team  string
	token string
}

func (s *slack) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/slack/message", func(w http.ResponseWriter, r *http.Request) {
		s.robot.Receive(&message{
			userId:      r.PostFormValue("user_id"),
			userName:    r.PostFormValue("user_name"),
			channelId:   r.PostFormValue("channel_id"),
			channelName: r.PostFormValue("channel_name"),
			text:        r.PostFormValue("text"),
		})
	}).Methods("POST")
}

func (s *slack) Send(channelId, msg string) {
	body, err := json.Marshal(&outgoingMessage{
		Channel:  channelId,
		Username: s.robot.Name(),
		Text:     msg,
	})

	if err != nil {
		log.Println("error sending to chat:", err)
	}

	endpoint := fmt.Sprintf("https://%s.slack.com/services/hooks/incoming-webhook?token=%s", s.team, s.token)
	log.Println(endpoint)
	log.Println(string(body))
	log.Println(http.PostForm(endpoint, url.Values{"payload": {string(body)}}))
}

func (s *slack) Stop() {
}

type outgoingMessage struct {
	Channel  string `json:"channel"`
	Username string `json:"username"`
	Text     string `json:"text"`
}

type message struct {
	userId      string
	userName    string
	channelId   string
	channelName string
	text        string
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
