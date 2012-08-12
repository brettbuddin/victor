package victor

import (
    "math/rand"
)

type TextMessage struct {
    Id        int
    Body      string
    CreatedAt string

    matches []string

    Send  func(text string)
    Reply func(text string)
    Paste func(text string)
    Sound func(text string)
}

func (self *TextMessage) SetMatches(matches []string) {
    self.matches = matches
}

func (self *TextMessage) Matches() []string {
    return self.matches
}

func (self *TextMessage) RandomString(strings []string) string {
    return strings[rand.Intn(len(strings))]
}
