package main

import (
    "github.com/brettbuddin/victor"
)

func main() {
    //r := victor.NewShell("victor")
    r := victor.NewCampfire("victor", "account", "token", []int{12345})

    r.Hear("derp", func(msg *victor.TextMessage) {
        msg.Send("Derp!")
    })

    r.Respond("hello", func(msg *victor.TextMessage) {
        msg.Reply("Hello!")
    })

    r.Run()
}
