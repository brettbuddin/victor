package victor

import (
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

type Robot struct {
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

// New returns a Robot
func New(adapterName, robotName, httpAddr string) *Robot {
	initFunc, err := chat.Load(adapterName)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	bot := &Robot{
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
func (r *Robot) Direct(exp string) string {
	return strings.Join([]string{
		"(?i)", // flags
		"\\A",  // begin
		"(?:(?:@)?" + r.name + "[:,]?\\s*|/)", // bot name
		"(?:" + exp + ")",                     // expression
		"\\z",                                 // end
	}, "")
}

func (r *Robot) Receive(m chat.Message) {
	r.incoming <- m
}

// Run starts the robot.
func (r *Robot) Run() error {
	go r.adapter.Run()
	go func() {
		for {
			select {
			case <-r.stop:
				return
			case m := <-r.incoming:
				if strings.ToLower(m.UserName()) != r.name {
					go r.Process(m)
				}
			}
		}
	}()

	r.http.Handle("/", r.router)
	return r.http.ListenAndServe(r.httpAddr)
}

func (r *Robot) Stop() {
	r.adapter.Stop()
	r.stop <- struct{}{}
	close(r.incoming)
	r.http.Stop()
}

func (r *Robot) Name() string {
	return r.name
}

func (r *Robot) SetName(n string) {
	r.name = n
}

func (r *Robot) Store() store.Store {
	return r.store
}

func (r *Robot) SetStore(s store.Store) {
	r.store = s
}

func (r *Robot) HTTP() *mux.Router {
	return r.router
}

func (r *Robot) SetHTTP(router *mux.Router) {
	r.router = router
}

func (r *Robot) Chat() chat.Adapter {
	return r.adapter
}

func (r *Robot) SetChat(name string) error {
	initFunc, err := chat.Load(name)

	if err != nil {
		return err
	}

	r.adapter = initFunc(r)
	return nil
}
