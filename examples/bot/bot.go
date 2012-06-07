/*
*   Usage:
*
*    $ export VICTOR_CAMPFIRE_ACCOUNT=
*    $ export VICTOR_CAMPFIRE_TOKEN=
*    $ export VICTOR_CAMPFIRE_ROOMS=
*    $ bot -adapter campfire -name ralph
*/

package main

import (
    "flag"
    "github.com/brettbuddin/victor"
)

func main() {
    adapter := flag.String("adapter", "", "victor adapter (campfire or shell)")
    name    := flag.String("name", "", "victor's new name in chat")

    flag.Parse()

    r := victor.NewRobot(*adapter, *name)

    r.Hear("derp", func(msg *victor.TextMessage) {
        msg.Send("Derp!")
    })

    r.Respond("hello", func(msg *victor.TextMessage) {
        msg.Reply("Hello!")
    })

    r.Run()
}
