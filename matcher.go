package victor

import (
    "regexp"
    "strings"
)

type Matcher struct {
    Exp      *regexp.Regexp
    Callback func(*TextMessage)
}

func NewMatcher(exp *regexp.Regexp, callback func(*TextMessage)) *Matcher {
    return &Matcher{
        Exp:      exp,
        Callback: callback,
    }
}

func (self *Matcher) Test(msg *TextMessage) bool {
    results := self.Exp.FindAllStringSubmatch(strings.ToLower(msg.Body), -1)

    if len(results) > 0 {
        msg.SetMatches(results[0])
        return true
    }

    return false
}
