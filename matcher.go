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

// Test matches the regexp pattern it's been told to apply to the incoming
// message and sets matches that result from it.
// Returns true if it is a match and false if it is not
func (self *Matcher) Test(msg *TextMessage) bool {
    results := self.Pattern.FindAllStringSubmatch(strings.ToLower(msg.Body), -1)

    if len(results) > 0 {
        msg.SetMatches(results[0])
        return true
    }

    return false
}
