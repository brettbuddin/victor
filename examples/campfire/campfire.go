package main

import (
    "github.com/brettbuddin/victor"
)

func main() {
    options := map[string]string{
        "account": "",          // name of the account (subdomain of *.campfirenow.com)
        "token": "",            // token for the user
        "rooms": "",            // comma seperated list
    }

    r := victor.NewRobot("campfire", options)

    r.Hear("derp", func(msg *victor.TextMessage) {
        msg.Send("Derp!")
    })

    r.Respond("hello", func(msg *victor.TextMessage) {
        msg.Reply("Hello!")
    })

    r.Run()
}
