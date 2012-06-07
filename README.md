**Victor** is a 37Signals Campfire bot written in Go (http://golang.org); inspired by Github's Hubot (http://github.com/github/hubot).

Here's a sample Victor executable:

```go
package main

import (
    "github.com/brettbuddin/victor"
)

func main() {
    r := victor.NewRobot("campfire", "victor")

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
