package victor

import (
    "regexp"
    "strings"
)

type Matcher struct {
    Pattern  *regexp.Regexp
    Callback func(*Context)
}

func NewMatcher(exp *regexp.Regexp, callback func(*Context)) *Matcher {
    return &Matcher{
        Pattern:  exp,
        Callback: callback,
    }
}

// Test matches the regexp pattern it's been told to apply to the incoming
// message and sets matches that result from it.
// Returns true if it is a match and false if it is not
func (m *Matcher) Test(ctx *Context) bool {
    results := m.Pattern.FindAllStringSubmatch(ctx.Message().Body(), -1)

    if len(results) > 0 {
        ctx.SetMatches(results[0])
        return true
    }

    return false
}
