/*
*   Usage:
*    $ bot -adapter campfire -name ralph
*/

package main

import (
    "flag"
    "github.com/brettbuddin/victor"
)

func main() {
    brain := victor.NewBrain("victor")
    r     := victor.NewCampfire(brain, "account", "token", [12345])
    //r   := victor.NewShell(brain)

    r.Hear("derp", func(msg *victor.TextMessage) {
        msg.Send("Derp!")
    })

    r.Respond("hello", func(msg *victor.TextMessage) {
        msg.Reply("Hello!")
    })

    r.Run()
}
