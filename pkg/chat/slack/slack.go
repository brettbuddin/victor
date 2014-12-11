package slack

import (
	"encoding/json"
	"fmt"
	"github.com/brettbuddin/victor/pkg/chat"
	"log"
	"net/http"
	"net/url"
	"os"
)

func init() {
	chat.Register("slack", func(r chat.Robot) chat.Adapter {
		team := os.Getenv("VICTOR_SLACK_TEAM")
		token := os.Getenv("VICTOR_SLACK_TOKEN")

		if team == "" || token == "" {
			log.Println("The following environment variables are required:")
			log.Println("VICTOR_SLACK_TEAM, VICTOR_SLACK_TOKEN")
			os.Exit(1)
		}

		return &slack{
			robot: r,
			team:  team,
			token: token,
		}
	})
}

type slack struct {
	robot       chat.Robot
	team, token string
}

func (s *slack) Run() {
	s.robot.HTTP().HandleFunc("/hubot/slack-webhook", func(w http.ResponseWriter, r *http.Request) {
		s.robot.Receive(&message{
			userID:      r.PostFormValue("user_id"),
			userName:    r.PostFormValue("user_name"),
			channelID:   r.PostFormValue("channel_id"),
			channelName: r.PostFormValue("channel_name"),
			text:        r.PostFormValue("text"),
		})
	}).Methods("POST")
}

func (s *slack) Send(channelID, msg string) {
	body, err := json.Marshal(&outgoingMessage{
		Channel:  channelID,
		Username: s.robot.Name(),
		Text:     msg,
	})

	if err != nil {
		log.Println("error sending to chat:", err)
	}

	endpoint := fmt.Sprintf("https://%s.slack.com/services/hooks/hubot?token=%s", s.team, s.token)
	http.PostForm(endpoint, url.Values{"payload": {string(body)}})
}

func (s *slack) Stop() {
}

type outgoingMessage struct {
	Channel  string `json:"channel"`
	Username string `json:"username"`
	Text     string `json:"text"`
}

type message struct {
	userID, userName, channelID, channelName, text string
}

func (m *message) UserID() string {
	return m.userID
}

func (m *message) UserName() string {
	return m.userName
}

func (m *message) ChannelID() string {
	return m.channelID
}

func (m *message) ChannelName() string {
	return m.channelName
}

func (m *message) Text() string {
	return m.text
}
