**Victor** is a 37Signals Campfire bot written in Go (http://golang.org); inspired by Github's Hubot (http://github.com/github/hubot).

Here's a sample Gobias executable:

```
package main

import (
    "victor"
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
```

### Robot Abilities

- `Hear`: Trigger an action based on some criteria heard anywhere in the channel.
- `Respond`: Respond to a direct statement at the robot (e.g. "gobias show me the diff")

### Actions on Messages

- `Send`: Send a bit of text to the channel.
- `Reply`: Reply directly to the person that triggered the action (e.g. "Brett: Yo yo yo. Here's the diff:")

### About

I'd been interested in Go for some time now so this was written an exercise to play with the language a bit. What I found while writing this is that I really like Go. A lot.
