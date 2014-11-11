package victor

import (
	"fmt"
	"github.com/brettbuddin/victor/pkg/chat"
	_ "github.com/brettbuddin/victor/pkg/chat/campfire"
	_ "github.com/brettbuddin/victor/pkg/chat/hipchat"
	_ "github.com/brettbuddin/victor/pkg/chat/shell"
	_ "github.com/brettbuddin/victor/pkg/chat/slack"
	"github.com/brettbuddin/victor/pkg/httpserver"
	"github.com/brettbuddin/victor/pkg/store"
	"github.com/gorilla/mux"
	"log"
	"os"
	"strings"
)

type Robot interface {
	Run()
	Stop()
	Name() string
	HandleFunc(string, HandlerFunc)
	Handle(string, Handler)
	Direct(string) string
	Receive(chat.Message)
	Chat() chat.Adapter
	Store() store.Adapter
	HTTP() *mux.Router
}

type Config struct {
	Name,
	ChatAdapter,
	StoreAdapter,
	HTTPAddr string
}

type robot struct {
	*dispatch
	name       string
	http       *httpserver.Server
	httpAddr   string
	httpRouter *mux.Router
	store      store.Adapter
	chat       chat.Adapter
	incoming   chan chat.Message
	stop       chan struct{}
}

// New returns a robot
func New(config Config) *robot {
	chatAdapter := config.ChatAdapter
	if chatAdapter == "" {
		chatAdapter = "shell"
	}

	chatInitFunc, err := chat.Load(config.ChatAdapter)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	storeAdapter := config.StoreAdapter
	if storeAdapter == "" {
		storeAdapter = "memory"
	}

	storeInitFunc, err := store.Load(storeAdapter)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	botName := config.Name
	if botName == "" {
		botName = "victor"
	}

	httpAddr := config.HTTPAddr
	if httpAddr == "" {
		httpAddr = ":9000"
	}

	bot := &robot{
		name:     botName,
		http:     httpserver.New(),
		httpAddr: httpAddr,
		incoming: make(chan chat.Message),
		store:    storeInitFunc(),
		stop:     make(chan struct{}),
	}

	bot.dispatch = newDispatch(bot)
	bot.chat = chatInitFunc(bot)
	bot.httpRouter = handlers(bot)

	defaults(bot)
	return bot
}

// Receive accepts messages for processing
func (r *robot) Receive(m chat.Message) {
	r.incoming <- m
}

// Run starts the robot.
func (r *robot) Run() {
	go r.chat.Run()
	go func() {
		for {
			select {
			case <-r.stop:
				close(r.incoming)
				return
			case m := <-r.incoming:
				if strings.ToLower(m.UserName()) != r.name {
					go r.ProcessMessage(m)
				}
			}
		}
	}()

	r.http.Handle("/", r.httpRouter)
	r.http.ListenAndServe(r.httpAddr)
}

// Stop shuts down the bot
func (r *robot) Stop() {
	r.chat.Stop()
	close(r.stop)
	r.http.Stop()
}

// Name returns the name of the bot
func (r *robot) Name() string {
	return r.name
}

// Store returns the data store adapter
func (r *robot) Store() store.Adapter {
	return r.store
}

// HTTP returns the HTTP router
func (r *robot) HTTP() *mux.Router {
	return r.httpRouter
}

// Chat returns the chat adapter
func (r *robot) Chat() chat.Adapter {
	return r.chat
}

// OnlyAllow provides a way of permitting specific users
// to execute a handler registered with the bot
func OnlyAllow(userNames []string, action func(s State)) func(State) {
	return func(s State) {
		actual := s.Message().UserName()
		for _, name := range userNames {
			if name == actual {
				action(s)
				return
			}
		}

		s.Chat().Send(s.Message().ChannelID(), fmt.Sprintf("Sorry, %s. I can't let you do that.", actual))
	}
}
