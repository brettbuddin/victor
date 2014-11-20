// Based on the Slack library for victor https://github.com/brettbuddin/victor

package hipchat

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/abachman/hipchat-go/hipchat"
	"github.com/brettbuddin/victor/pkg/chat"
)

func init() {
	chat.Register("hipchat", func(r chat.Robot) chat.Adapter {
		rooms := os.Getenv("VICTOR_HIPCHAT_ROOMS")
		token := os.Getenv("VICTOR_HIPCHAT_TOKEN")
		webhookUrl := os.Getenv("VICTOR_HIPCHAT_WEBHOOK")

		if rooms == "" || token == "" {
			log.Println("The following environment variables are required:")
			log.Println("VICTOR_HIPCHAT_ROOMS, VICTOR_HIPCHAT_TOKEN, VICTOR_HIPCHAT_WEBHOOK")
			os.Exit(1)
		}

		// A HipChat API v2 Client
		client := hipchat.NewClient(token)

		hook, _ := url.Parse(webhookUrl)

		hookName := "victor webhook"
		hookEvent := "room_message"
		hookUrl := fmt.Sprintf("%s://%s%s", hook.Scheme, hook.Host, hook.Path)

		// Initialize HipChat Webhooks - your bot must have admin access to perform
		// this action. Failure to create webhooks won't stop the bot from
		// launching, but you'll have to create them by some other means.
		var exists bool
		for _, roomId := range strings.Split(rooms, ",") {
			hooks, resp, err := client.Room.GetAllWebhooks(roomId, nil)
			handleRequestError(resp, err)
			if err != nil {
				continue
			}

			exists = false
			for _, webhook := range hooks.Webhooks {
				if webhook.Event == hookEvent && webhook.URL == hookUrl {
					exists = true
					break
				}
			}

			if !exists {
				_, resp, err := client.Room.CreateWebhook(roomId, &hipchat.CreateWebhookRequest{
					Name:  hookName,
					Event: hookEvent,
					URL:   hookUrl,
				})
				handleRequestError(resp, err)

				if err != nil {
					log.Println("[hipchat adapter init] failed to create webhook")
				}
			}
		}

		return &adapter{
			robot:  r,
			client: client,
			hook:   hook,
		}
	})
}

type adapter struct {
	robot  chat.Robot
	client *hipchat.Client
	hook   *url.URL
	rooms  []*hipchat.Room
}

func (h *adapter) Run() {
	h.robot.HTTP().HandleFunc(h.hook.Path, func(w http.ResponseWriter, r *http.Request) {
		// debug message JSON
		body, _ := ioutil.ReadAll(r.Body)

		msg := &WebhookMessage{}
		err := json.NewDecoder(strings.NewReader(string(body))).Decode(msg)
		if err != nil {
			log.Println("[hipchat adapter Run] failed to decode message:", err)
			log.Println(string(body))
		}

		// now communicate message to victor
		h.robot.Receive(&message{
			userID:      strconv.Itoa(msg.Item.Message.From.ID),
			userName:    msg.Item.Message.From.Name,
			channelID:   strconv.Itoa(msg.Item.Room.ID),
			channelName: msg.Item.Room.Name,
			text:        msg.Item.Message.Message,
		})
	}).Methods("POST")
}

// additional option with Hipchat adapter since Hipchat allows specifying HTML
func (h *adapter) SendHtml(channelID, msg string) {
	resp, err := h.client.Room.Notification(channelID, &hipchat.NotificationRequest{
		Message:       msg,
		Notify:        true,
		MessageFormat: "html",
	})

	handleRequestError(resp, err)
}

func (h *adapter) Send(channelID, msg string) {
	resp, err := h.client.Room.Notification(channelID, &hipchat.NotificationRequest{
		Message:       msg,
		Notify:        true,
		MessageFormat: "text",
	})

	handleRequestError(resp, err)
}

func (s *adapter) Stop() {
	// no need, webhook listeners are stopped with the HTTP server
}

// victor boilerplate
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

// report API errors
func handleRequestError(resp *http.Response, err error) {
	if err != nil {
		if resp != nil {
			log.Printf("[hipchat adapter] request failed: %+v\n", resp)
			body, _ := ioutil.ReadAll(resp.Body)
			log.Printf("%+v\n", body)
		} else {
			log.Println("[hipchat adapter] request failed, response is nil")
		}
		log.Println(err)
	}
}

// hipchat webhook API
type User struct {
	ID          int    `json:"id"`
	MentionName string `json:"mention_name"`
	Name        string `json:"name"`
}

type Message struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date,omitempty"`
	Message string `json:"message,omitempty"`
	From    User   `json:"from,omitempty"`
}

type Room struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MessageData struct {
	Message Message `json:"message"`
	Room    Room    `json:"room"`
}

// top level hipchat webhook JSON object
type WebhookMessage struct {
	Event         string      `json:"string"`
	Item          MessageData `json:"item"`
	WebhookID     int         `json:"webhook_id"`
	OauthClientID string      `json:"oauth_client_id"`
}
