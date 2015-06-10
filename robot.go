package victor

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/brettbuddin/victor/pkg/chat"
	_ "github.com/brettbuddin/victor/pkg/chat/campfire"
	_ "github.com/brettbuddin/victor/pkg/chat/hipchat"
	_ "github.com/brettbuddin/victor/pkg/chat/shell"
	_ "github.com/brettbuddin/victor/pkg/chat/slack"
	_ "github.com/brettbuddin/victor/pkg/chat/slackRealtime"
	"github.com/brettbuddin/victor/pkg/httpserver"
	"github.com/brettbuddin/victor/pkg/store"
	_ "github.com/brettbuddin/victor/pkg/store/boltstore"
	_ "github.com/brettbuddin/victor/pkg/store/memory"
	"github.com/gorilla/mux"
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
	AdapterConfig() (interface{}, bool)
	StoreConfig() (interface{}, bool)
}

type Config struct {
	Name,
	ChatAdapter,
	StoreAdapter,
	HTTPAddr string
	AdapterConfig,
	StoreConfig interface{}
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
	adapterConfig,
	storeConfig interface{}
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
		httpAddr: httpAddr,
		incoming: make(chan chat.Message),
		stop:     make(chan struct{}),
	}

	bot.store = storeInitFunc(bot)
	bot.adapterConfig = config.AdapterConfig
	bot.dispatch = newDispatch(bot)
	bot.chat = chatInitFunc(bot)
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

// HTTP returns the HTTP router.
// The HTTP router is disabled (uninitialized) by default but is created upon
// the first call to HTTP().
//
// TODO consider having one explicitly enable the storage access endpoints
// (defined in http_handlers.go) since, as far as I can tell, a chat/storage
// adapter might want to use the included http router without enabling access
// to the entire store via those endpoints. At the moment these are coupled
// together.
func (r *robot) HTTP() *mux.Router {
	if r.httpRouter == nil {
		r.initHTTP()
	}
	return r.httpRouter
}

// Chat returns the chat adapter
func (r *robot) Chat() chat.Adapter {
	return r.chat
}

func (r *robot) AdapterConfig() (interface{}, bool) {
	return r.adapterConfig, r.adapterConfig != nil
}

func (r *robot) StoreConfig() (interface{}, bool) {
	return r.storeConfig, r.storeConfig != nil
}

func (r *robot) initHTTP() {
	log.Println("Initializing victor's HTTP server.")
	r.http = httpserver.New()
	r.httpRouter = handlers(r)
	r.http.Handle("/", r.httpRouter)
	r.http.ListenAndServe(r.httpAddr)
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
