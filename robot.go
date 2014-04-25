package victor

import (
	"fmt"
	"github.com/brettbuddin/victor/pkg/chat"
	_ "github.com/brettbuddin/victor/pkg/chat/campfire"
	_ "github.com/brettbuddin/victor/pkg/chat/shell"
	_ "github.com/brettbuddin/victor/pkg/chat/slack"
	"github.com/brettbuddin/victor/pkg/httpserver"
	"github.com/brettbuddin/victor/pkg/store"
	"github.com/gorilla/mux"
	"log"
	"os"
	"strings"
)

type Runner interface {
	Run() error
	Stop()
}

type Chatter interface {
	Name() string
	HandleFunc(string, HandlerFunc)
	Handle(string, Handler)
	Direct(string) string
	Receive(chat.Message)
	Chat() chat.Adapter
}

type Persister interface {
	Store() store.Adapter
}

type Robot interface {
	Runner
	Chatter
	Persister

	HTTP() *mux.Router
}

type Config struct {
	Name         string
	ChatAdapter  string
	StoreAdapter string
	HTTPAddr     string
}

type robot struct {
	*dispatch
	name     string
	http     *httpserver.Server
	httpAddr string
	router   *mux.Router
	store    store.Adapter
	chat     chat.Adapter
	incoming chan chat.Message
	stop     chan struct{}
}

// New returns a robot
func New(config Config) *robot {
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

	bot := &robot{
		name:     config.Name,
		http:     httpserver.New(),
		httpAddr: config.HTTPAddr,
		incoming: make(chan chat.Message),
		store:    storeInitFunc(),
		stop:     make(chan struct{}),
	}

	bot.dispatch = newDispatch(bot)
	bot.chat = chatInitFunc(bot)
	bot.router = handlers(bot)

	defaults(bot)
	return bot
}

// Direct wraps a regexp pattern in the necessary pattern
// for a direct command to the bot.
func (r *robot) Direct(exp string) string {
	return strings.Join([]string{
		"(?i)", // flags
		"\\A",  // begin
		"(?:(?:@)?" + r.name + "[:,]?\\s*|/)", // bot name
		"(?:" + exp + ")",                     // expression
		"\\z",                                 // end
	}, "")
}

func (r *robot) Receive(m chat.Message) {
	r.incoming <- m
}

// Run starts the robot.
func (r *robot) Run() error {
	go r.chat.Run()
	go func() {
		for {
			select {
			case <-r.stop:
				return
			case m := <-r.incoming:
				if strings.ToLower(m.UserName()) != r.name {
					go r.ProcessMessage(m)
				}
			}
		}
	}()

	r.http.Handle("/", r.router)
	return r.http.ListenAndServe(r.httpAddr)
}

func (r *robot) Stop() {
	r.chat.Stop()
	r.stop <- struct{}{}
	close(r.incoming)
	r.http.Stop()
}

func (r *robot) Name() string {
	return r.name
}

func (r *robot) Store() store.Adapter {
	return r.store
}

func (r *robot) HTTP() *mux.Router {
	return r.router
}

func (r *robot) Chat() chat.Adapter {
	return r.chat
}

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
