package victor

import (
    "log"
    "regexp"
    "strings"

    "github.com/brettbuddin/victor/utils/google"
)

type RobotAdapter interface {
    Hear(string, func(*TextMessage))
    Respond(string, func(*TextMessage))
    Run()
}

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

    brain.registerDefaultAbilities()

    return brain
}

func (self *Brain) Hear(expStr string, callback func(*TextMessage)) {
    exp, _ := regexp.Compile(expStr)
    self.matchers = append(self.matchers, NewMatcher(exp, callback))

    log.Printf("Pattern: /%s/", exp.String())
}

func (self *Brain) Respond(expStr string, callback func(*TextMessage)) {
    expWithNameStr := "^(" + self.name + "[:,]?)\\s*(?:" + expStr + ")"
    exp, _ := regexp.Compile(strings.ToLower(expWithNameStr))

    self.matchers = append(self.matchers, NewMatcher(exp, callback))

    log.Printf("Pattern: /%s/", exp.String())
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

func (self *Brain) registerDefaultAbilities() {
    self.Respond("ping", func(msg *TextMessage) {
        msg.Send("Pong!")
    })

    self.Respond("(image|img) (.*)", func(msg *TextMessage) {
        result, err := google.ImageSearch(msg.Matches()[3])

        if err != nil {
            log.Print(err)
            return
        }

        if result == "" {
            msg.Send("I didn't find anything.")
            return
        }

        msg.Send(result)
    })

    self.Respond("(list|show) users", func(msg *TextMessage) {
        list := ""

        for _, user := range self.users {
            list += user.Name + "\n"
        }

        msg.Paste(list)
    })
}
