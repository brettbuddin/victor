package main

import (
    "robot"
)

func main() {
    options := map[string]string{
        "name": "henry",
    }

    r := victor.NewRobot("shell", options)

    r.Hear("derp", func(msg *victor.TextMessage) {
        msg.Send("Derp!")
    })

    r.Respond("hello", func(msg *victor.TextMessage) {
        msg.Reply("Hello!")
    })

    r.Run()
}
