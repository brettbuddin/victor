package victor

import (
    "regexp"
    "strings"
)

type Matcher struct {
    Pattern  *regexp.Regexp
    Callback func(*TextMessage)
}

func NewMatcher(exp *regexp.Regexp, callback func(*TextMessage)) *Matcher {
    return &Matcher{
        Pattern:  exp,
        Callback: callback,
    }
}

func (self *Matcher) Test(msg *TextMessage) bool {
    results := self.Pattern.FindAllStringSubmatch(strings.ToLower(msg.Body), -1)

    if len(results) > 0 {
        msg.SetMatches(results[0])
        return true
    }

    return false
}
