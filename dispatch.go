package victor

import (
	"github.com/brettbuddin/victor/pkg/chat"
	"regexp"
	"strings"
)

type dispatch struct {
	robot    Robot
	handlers map[*regexp.Regexp]Handler
}

func newDispatch(bot Robot) *dispatch {
	return &dispatch{
		robot:    bot,
		handlers: make(map[*regexp.Regexp]Handler),
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
	d.handlers[regexp.MustCompile(exp)] = h
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
	for exp, handler := range d.handlers {
		matches := exp.FindAllStringSubmatch(m.Text(), -1)

		if len(matches) > 0 {
			params := matches[0][1:]
			handler.Handle(&state{
				robot:   d.robot,
				message: m,
				params:  params,
			})
			return
		}
	}
}
