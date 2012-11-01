package victor

import (
    "math/rand"
)

type Message struct {
    Id        int
    Type      string
    UserId    int
    RoomId    int
    Body      string
    CreatedAt string
}

type Context struct {
    message *Message

    Send  func(text string)
    Reply func(text string)
    Paste func(text string)
    Sound func(text string)

    matches []string
}

func (self *Context) SetMessage(msg *Message) {
    self.message = msg
}

func (self *Context) Message() *Message {
    return self.message
}

func (self *Context) SetMatches(matches []string) {
    self.matches = matches
}

func (self *Context) Matches() []string {
    return self.matches
}

func (self *Context) RandomString(strings []string) string {
    return strings[rand.Intn(len(strings))]
}
