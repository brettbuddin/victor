package victor

import (
	"github.com/brettbuddin/victor/pkg/chat"
	"regexp"
)

type Dispatch struct {
	robot    Robot
	handlers map[*regexp.Regexp]Handler
}

func NewDispatch(bot Robot) *Dispatch {
	return &Dispatch{
		robot:    bot,
		handlers: make(map[*regexp.Regexp]Handler),
	}
}

func (d *Dispatch) Handle(exp string, h Handler) {
	d.handlers[regexp.MustCompile(exp)] = h
}

func (d *Dispatch) HandleFunc(exp string, f HandlerFunc) {
	d.Handle(exp, f)
}

func (d *Dispatch) ProcessMessage(m chat.Message) {
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
