package victor_test

import (
    "github.com/brettbuddin/victor"
    "testing"
)

func TestHear(t *testing.T) {
    doesHear(t, "d[eu]rp", "derp")
    doesHear(t, "d[eu]rp", "durp")

    doesNotHear(t, "hello$", "hello nope")
    doesNotHear(t, "^hello", "nope hello")
}

func TestRespond(t *testing.T) {
    doesRespond(t, "d[eu]rp", "robot: derp")
    doesRespond(t, "d[eu]rp", "robot durp")
    doesRespond(t, "d[eu]rp", "robot, durp")
    doesRespond(t, "d[eu]rp$", "robot, durp")
    doesRespond(t, "d[eu]rp$", "Robot, durp")
}

func doesRespond(t *testing.T, pattern, messageText string) {
    if testRespond(pattern, messageText) == false {
        t.Errorf("%s not triggered by %s", pattern, messageText)
    }
}

func doesNotRespond(t *testing.T, pattern, messageText string) {
    if testRespond(pattern, messageText) {
        t.Errorf("%s triggered by %s", pattern, messageText)
    }
}

func doesHear(t *testing.T, pattern, messageText string) {
    if testHear(pattern, messageText) == false {
        t.Errorf("%s not triggered by %s", pattern, messageText)
    }
}

func doesNotHear(t *testing.T, pattern, messageText string) {
    if testHear(pattern, messageText) {
        t.Errorf("%s triggered by %s", pattern, messageText)
    }
}

func testRespond(pattern, messageText string) bool {
    brain := victor.NewBrain("robot")
    triggered := false

    brain.Respond(pattern, func(msg *victor.TextMessage) {
        triggered = true
    })

    msg := &victor.TextMessage{
        Id:   1,
        Body: messageText,
    }

    brain.Receive(msg)

    return triggered
}

func testHear(pattern, messageText string) bool {
    brain := victor.NewBrain("robot")
    triggered := false

    brain.Hear(pattern, func(msg *victor.TextMessage) {
        triggered = true
    })

    msg := &victor.TextMessage{
        Id:   1,
        Body: messageText,
    }

    brain.Receive(msg)

    return triggered
}
