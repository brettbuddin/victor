## Ghost

**Ghost** is a library for creating your own chat bot.

We use Ghost as the backbone of our bot, Virbot, within our team Campfire at Virb (http://virb.com).

### Making Him Your Own

Ghost is more of a framework for constructing your own bot so he doesn't do a whole lot out-of-the-box. I'll be adding more default behavior to him as time progresses, but you might want him to do something specific to your team's needs. You can use the programs located in `examples/` as starting points to create your own executable.

### Listening for Things

There are two ways to trigger actions on the bot:

- `Hear`: Trigger an action based on some criteria heard anywhere in the channel.
- `Respond`: Respond to a direct statement at the robot (e.g. "virbot show not shipped")

```go
bot := ghost.New("campfire", "ralph")

// Triggers anytime it hears the word "alot".
bot.Hear("alot", func(m robot.Message) {
    m.Room().Say("A LOT.")
})

// Triggers when someone talks directly to the bot: "ralph hi"
bot.Respond("hi|hello|howdy", func(m robot.Message) {
    m.Reply(fmt.Sprintf("Hello, %s", m.User().Name()))
})

bot.Respond("bye", func(m robot.Message) {
    m.Reply(fmt.Sprintf("Goodbye, %s", m.User().Name()))
})

bot.Run()
```
