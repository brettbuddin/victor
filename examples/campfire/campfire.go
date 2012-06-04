package main

import (
    "robot"
)

func main() {
    options := map[string]string{
        "account": "",          // name of the account (subdomain of *.campfirenow.com)
        "token": "",            // token for the user
        "rooms": "",            // comma seperated list
    }

    r := robot.NewRobot("campfire", options)

    r.Hear("derp", func(msg *robot.TextMessage) {
        msg.Send("Derp!")
    })

    r.Respond("hello", func(msg *robot.TextMessage) {
        msg.Reply("Hello!")
    })

    r.Run()
}
