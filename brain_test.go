package victor_test

import (
    "github.com/brettbuddin/victor"
    "github.com/brettbuddin/victor/campfire"
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

    brain.Respond(pattern, func(ctx *victor.Context) {
        triggered = true
    })

    ctx := new(victor.Context)

    msg := &campfire.Message{}
    msg.SetId(1)
    msg.SetBody(messageText)

    ctx.SetMessage(msg)

    brain.Receive(ctx)

    return triggered
}

func testHear(pattern, messageText string) bool {
    brain := victor.NewBrain("robot")
    triggered := false

    brain.Hear(pattern, func(msg *victor.Context) {
        triggered = true
    })

    ctx := new(victor.Context)

    msg := &campfire.Message{}
    msg.SetId(1)
    msg.SetBody(messageText)

    ctx.SetMessage(msg)

    brain.Receive(ctx)

    return triggered
}
