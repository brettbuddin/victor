package victor

import (
    "regexp"
)

type Listener struct {
    Exp      *regexp.Regexp
    Callback func(*TextMessage)
}

func NewListener(exp *regexp.Regexp, callback func(*TextMessage)) *Listener {
    return &Listener{
        Exp: exp,
        Callback: callback,
    }
}

func (self *Listener) Test(msg *TextMessage) bool {
    return self.Exp.Match([]byte(msg.Body))
}
