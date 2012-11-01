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

func (self *Context) SetMessage(msg *Message) *Context {
    self.message = msg

    return self
}

func (self *Context) Message() *Message {
    return self.message
}

func (self *Context) SetMatches(matches []string) *Context {
    self.matches = matches

    return self
}

func (self *Context) Matches() []string {
    return self.matches
}

func (self *Context) RandomString(strings []string) string {
    return strings[rand.Intn(len(strings))]
}
