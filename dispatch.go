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

// Handle registers a Handler for matching
func (d *dispatch) Handle(exp string, h Handler) {
	d.handlers[regexp.MustCompile(exp)] = h
}

// HandleFunc registers a HandlerFunc for matching
func (d *dispatch) HandleFunc(exp string, f HandlerFunc) {
	d.Handle(exp, f)
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
