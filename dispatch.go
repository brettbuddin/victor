package victor

import (
	"regexp"
	"strings"

	"github.com/brettbuddin/victor/pkg/chat"
)

// HandlerPair provides an interface for a handler as well as the regular
// expression which a message should match in order to pass control onto the
// handler
type HandlerPair interface {
	Exp() *regexp.Regexp
	Handler() Handler
}

type handlerPair struct {
	exp    *regexp.Regexp
	handle Handler
}

func (pair *handlerPair) Exp() *regexp.Regexp {
	return pair.exp
}

func (pair *handlerPair) Handler() Handler {
	return pair.handle
}

type dispatch struct {
	robot    Robot
	handlers []HandlerPair
}

func newDispatch(bot Robot) *dispatch {
	return &dispatch{
		robot:    bot,
		handlers: make([]HandlerPair, 0, 10),
	}
}

// HandleCommand registers a Handler for matching statements directed at the bot
func (d *dispatch) HandleCommand(exp string, h Handler) {
	d.handle(d.Direct(exp), h)
}

// HandleCommandFunc registers a Handler for matching statements directed at the bot
func (d *dispatch) HandleCommandFunc(exp string, f HandlerFunc) {
	d.handle(d.Direct(exp), f)
}

// Handle registers a Handler for matching
func (d *dispatch) Handle(exp string, h Handler) {
	d.handle(exp, h)
}

// HandleFunc registers a HandlerFunc for matching
func (d *dispatch) HandleFunc(exp string, f HandlerFunc) {
	d.handle(exp, f)
}

func (d *dispatch) handle(exp string, h Handler) {
	d.handlers = append(d.handlers, &handlerPair{
		exp:    regexp.MustCompile(exp),
		handle: h,
	})
}

// Direct wraps a regexp pattern in the necessary pattern
// for a direct command to the bot.
func (d *dispatch) Direct(exp string) string {
	return strings.Join([]string{
		"(?i)", // flags
		"\\A",  // begin
		"(?:(?:@)?" + d.robot.Name() + "[:,]?\\s*|/)", // bot name
		"(?:" + exp + ")",                             // expression
		"\\z",                                         // end
	}, "")
}

// ProcessMessage finds a match for a message and runs its Handler
func (d *dispatch) ProcessMessage(m chat.Message) {
	for _, pair := range d.handlers {
		matches := pair.Exp().FindAllStringSubmatch(m.Text(), -1)

		if len(matches) > 0 {
			params := matches[0][1:]
			pair.Handler().Handle(&state{
				robot:   d.robot,
				message: m,
				params:  params,
			})
			return
		}
	}
}
