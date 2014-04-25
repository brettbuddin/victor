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

type Robot interface {
	HandleFunc(string, HandlerFunc)
	Handle(string, Handler)
	Direct(string) string
	Receive(chat.Message)
	Chat() chat.Adapter
	Store() store.Store
}

type robot struct {
	*Dispatch
	name     string
	http     *httpserver.Server
	httpAddr string
	router   *mux.Router
	store    store.Store
	adapter  chat.Adapter
	incoming chan chat.Message
	stop     chan struct{}
}

// New returns a robot
func New(adapterName, robotName, httpAddr string) *robot {
	initFunc, err := chat.Load(adapterName)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	bot := &robot{
		name:     robotName,
		http:     httpserver.New(),
		store:    store.NewMemoryStore(),
		incoming: make(chan chat.Message),
		httpAddr: httpAddr,
		stop:     make(chan struct{}),
	}

	bot.Dispatch = NewDispatch(bot)
	bot.adapter = initFunc(bot)
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
	go r.adapter.Run()
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
	r.adapter.Stop()
	r.stop <- struct{}{}
	close(r.incoming)
	r.http.Stop()
}

func (r *robot) Name() string {
	return r.name
}

func (r *robot) SetName(n string) {
	r.name = n
}

func (r *robot) Store() store.Store {
	return r.store
}

func (r *robot) SetStore(s store.Store) {
	r.store = s
}

func (r *robot) HTTP() *mux.Router {
	return r.router
}

func (r *robot) SetHTTP(router *mux.Router) {
	r.router = router
}

func (r *robot) Chat() chat.Adapter {
	return r.adapter
}

func (r *robot) SetChat(name string) error {
	initFunc, err := chat.Load(name)

	if err != nil {
		return err
	}

	r.adapter = initFunc(r)
	return nil
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
