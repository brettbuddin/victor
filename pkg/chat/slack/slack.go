package slack

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/brettbuddin/victor/pkg/chat"
)

func init() {
	chat.Register("slack", func(r chat.Robot) chat.Adapter {
		slackOutgoingWebhookPath := os.Getenv("SLACK_OUTGOING_WEBHOOK")
		slackIncomingWebhookUri := os.Getenv("SLACK_INCOMING_WEBHOOK_URI")

		if slackOutgoingWebhookPath == "" || slackIncomingWebhookUri == "" {
			log.Println("The following environment variable is required:")
			log.Println("SLACK_INCOMING_WEBHOOK_URI, SLACK_OUTGOING_WEBHOOK")
			os.Exit(1)
		}

		return &slack{
			robot:               r,
			outgoingWebhookPath: slackOutgoingWebhookPath,
			incomingWebhookUri:  slackIncomingWebhookUri,
		}
	})
}

type slack struct {
	robot               chat.Robot
	incomingWebhookUri  string
	outgoingWebhookPath string
}

func (s *slack) Run() {
	s.robot.HTTP().HandleFunc(s.outgoingWebhookPath, func(w http.ResponseWriter, r *http.Request) {
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

	if resp, err := HttpClient.PostForm(s.incomingWebhookUri, url.Values{"payload": {string(body)}}); err != nil {
		log.Printf("Slack API call failed: %s", err)
	} else if resp.StatusCode != 200 {
		log.Printf("Slack API returned HTTP status code %d", resp.StatusCode)
	}
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
