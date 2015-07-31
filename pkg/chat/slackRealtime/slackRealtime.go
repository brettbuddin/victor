package slackRealtime

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/brettbuddin/victor/pkg/chat"
	"github.com/nlopes/slack"
)

// The Slack Websocket's registered adapter name for the victor framework.
const AdapterName = "slackRealtime"

// Prefix for the user's ID which is used when reading/writing from the bot's store
const userInfoPrefix = AdapterName + "."

// init registers SlackAdapter to the victor chat framework.
func init() {
	chat.Register(AdapterName, func(r chat.Robot) chat.Adapter {
		config, configSet := r.AdapterConfig()
		if !configSet {
			log.Println("A configuration struct implementing the SlackConfig interface must be set.")
			os.Exit(1)
		}
		sConfig, ok := config.(Config)
		if !ok {
			log.Println("The bot's config must implement the SlackConfig interface.")
			os.Exit(1)
		}
		return &SlackAdapter{
			robot:      r,
			chSender:   make(chan *slack.OutgoingMessage),
			chReceiver: make(chan slack.SlackEvent),
			token:      sConfig.Token(),
		}
	})
}

// Config provides the slack adapter with the necessary
// information to open a websocket connection with the slack Real time API.
type Config interface {
	Token() string
}

// Config implements the SlackRealtimeConfig interface to provide a slack
// adapter with the information it needs to authenticate with slack.
type configImpl struct {
	token string
}

// NewConfig returns a new slack configuration instance using the given token.
func NewConfig(token string) configImpl {
	return configImpl{token: token}
}

// Token returns the slack token.
func (c configImpl) Token() string {
	return c.token
}

// SlackAdapter holds all information needed by the adapter to send/receive messages.
type SlackAdapter struct {
	robot      chat.Robot
	token      string
	instance   *slack.Slack
	wsAPI      *slack.WS
	chSender   chan *slack.OutgoingMessage
	chReceiver chan slack.SlackEvent
}

// Run starts the adapter and begins to listen for new messages to send/receive.
// At the moment this will crash the program and print the error messages to a
// log if the connection fails.
func (adapter *SlackAdapter) Run() {
	adapter.instance = slack.New(adapter.token)
	adapter.instance.SetDebug(false)
	// TODO need to look up what these values actually mean...
	var err error
	adapter.wsAPI, err = adapter.instance.StartRTM("", "http://example.com")
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	// sets up the monitoring code for sending/receiving messages from slack
	go adapter.wsAPI.HandleIncomingEvents(adapter.chReceiver)
	go adapter.wsAPI.Keepalive(20 * time.Second)
	adapter.monitorEvents()
}

// Stop stops the adapter.
// TODO implement
func (adapter *SlackAdapter) Stop() {
}

func (adapter *SlackAdapter) getUser(userID string) (*slack.User, error) {
	// try to get the stored user info
	userStr, exists := adapter.getStoreKey(userID)
	// if it hasn't been stored then perform a slack API call to get it and
	// store it
	if !exists {
		user, err := adapter.instance.GetUserInfo(userID)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		// try to encode it as a json string for storage
		var userArr []byte
		if userArr, err = json.Marshal(user); err == nil {
			adapter.setStoreKey(userID, string(userArr))
		}
		return user, nil
	}
	var user slack.User
	// convert the json string to the user object
	err := json.Unmarshal([]byte(userStr), &user)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &user, nil

	// TODO handle an error on the unmarshalling of the stored json object?
}

func (adapter *SlackAdapter) handleMessage(event *slack.MessageEvent) {
	user, _ := adapter.getUser(event.User)
	// TODO use error
	if user != nil {
		// ignore any messages that are sent by us
		if user.ID == adapter.instance.GetInfo().User.ID {
			return
		}
		msg := slackMessage{
			user:      user,
			text:      adapter.unescapeMessage(event.Text),
			channelID: event.Channel,
			// TODO change or not needed?
			channelName: event.Channel,
		}
		adapter.robot.Receive(&msg)
		log.Println(msg.Text())
	}
}

// Replace all instances of the bot's encoded name with it's actual name.
//
// TODO might want to update this to replace all matches of <@USER_ID> with
// the user's name.
func (adapter *SlackAdapter) unescapeMessage(msg string) string {
	userID := adapter.instance.GetInfo().User.ID
	return strings.Replace(msg, getEncodedUserID(userID), adapter.robot.Name(), -1)
}

// Returns the encoded string version of a user's slack ID.
func getEncodedUserID(userID string) string {
	return fmt.Sprintf("<@%s>", userID)
}

// monitorEvents handles incoming events and filters them to only worry about
// incoming messages.
func (adapter *SlackAdapter) monitorEvents() {
	for {
		msg := <-adapter.chReceiver
		switch msg.Data.(type) {
		case *slack.MessageEvent:
			go adapter.handleMessage(msg.Data.(*slack.MessageEvent))
		}
	}
}

// Send sends a message to the given slack channel.
func (adapter *SlackAdapter) Send(channelID, msg string) {
	msgObj := adapter.wsAPI.NewOutgoingMessage(msg, channelID)
	adapter.wsAPI.SendMessage(msgObj)
}

// getStoreKey is a helper method to access the robot's store.
func (adapter *SlackAdapter) getStoreKey(key string) (string, bool) {
	return adapter.robot.Store().Get(userInfoPrefix + key)
}

// setStoreKey is a helper method to access the robot's store.
func (adapter *SlackAdapter) setStoreKey(key, val string) {
	adapter.robot.Store().Set(userInfoPrefix+key, val)
}

// slackMessage is an internal struct implementing victor's message interface.
type slackMessage struct {
	user        *slack.User
	text        string
	channelID   string
	channelName string
}

func (m *slackMessage) User() *slack.User {
	return m.user
}

func (m *slackMessage) UserID() string {
	return m.user.ID
}

func (m *slackMessage) UserName() string {
	return m.user.Name
}

func (m *slackMessage) EmailAddress() string {
	return m.user.Profile.Email
}

func (m *slackMessage) ChannelID() string {
	return m.channelID
}

func (m *slackMessage) ChannelName() string {
	return m.channelName
}

func (m *slackMessage) Text() string {
	return m.text
}
