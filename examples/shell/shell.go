package main

import (
    "robot"
)

func main() {
    options := map[string]string{
        "name": "henry",
    }

    r := robot.NewRobot("shell", options)

    r.Hear("derp", func(msg *robot.TextMessage) {
        msg.Send("Derp!")
    })

    r.Respond("hello", func(msg *robot.TextMessage) {
        msg.Reply("Hello!")
    })

    r.Run()
}
