package victor

import (
    "regexp"
    "strings"
)

type Listener struct {
    Exp      *regexp.Regexp
    Callback func(*TextMessage)
}

func NewListener(exp *regexp.Regexp, callback func(*TextMessage)) *Listener {
    return &Listener{
        Exp:      exp,
        Callback: callback,
    }
}

func (self *Listener) Test(msg *TextMessage) bool {
    results := self.Exp.FindAllStringSubmatch(strings.ToLower(msg.Body), -1)

    if len(results) > 0 {
        msg.SetMatches(results[0])
        return true
    }

    return false
}
