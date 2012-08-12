package victor

import (
    "log"
    "regexp"
    "strings"
)

type Brain struct {
    name     string
    options  map[string]string
    matchers []*Matcher
    users    []*User
}

func NewBrain(name string) *Brain {
    brain := &Brain{
        name:     name,
        matchers: make([]*Matcher, 0, 1),
    }

    RegisterDefaultAbilities(brain)

    return brain
}

func (self *Brain) AddMatcher(m *Matcher) {
    self.matchers = append(self.matchers, m)

    log.Printf("Pattern: /%s/", m.Pattern.String())
}

func (self *Brain) Hear(expStr string, callback func(*TextMessage)) {
    exp, _ := regexp.Compile(expStr)

    self.AddMatcher(NewMatcher(exp, callback))
}

func (self *Brain) Respond(expStr string, callback func(*TextMessage)) {
    expWithNameStr := "^(" + self.name + "[:,]?)\\s*(?:" + expStr + ")"
    exp, _ := regexp.Compile(strings.ToLower(expWithNameStr))

    self.AddMatcher(NewMatcher(exp, callback))
}

func (self *Brain) Receive(msg *TextMessage) {
    for _, matcher := range self.matchers {
        if matcher.Test(msg) {
            matcher.Callback(msg)
        }
    }
}

func (self *Brain) RememberUser(user *User) {
    for i, u := range self.users {
        if u.Id == user.Id {
            // update the name if its different
            if u.Name != user.Name {
                self.users[i].Name = user.Name
            }

            return
        }
    }

    self.users = append(self.users, user)
}

func (self *Brain) Users() []*User {
    return self.users
}

func (self *Brain) UserForId(id int) *User {
    for _, user := range self.users {
        if user.Id == id {
            return user
        }
    }

    return nil
}

func (self *Brain) Shutdown() {
    log.Print("See ya!")
}
