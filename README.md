## Victor

[![Build Status](https://travis-ci.org/brettbuddin/victor.png?branch=master)](https://travis-ci.org/brettbuddin/victor)

**Victor** is a library for creating your own chat bot.

We use Victor as the backbone of our bot, Virbot, within our team Campfire at Virb (http://virb.com).

### Making Him Your Own

Victor is more of a framework for constructing your own bot so he doesn't do a whole lot out-of-the-box. I'll be adding more default behavior to him as time progresses, but you might want him to do something specific to your team's needs. You can use the programs located in `examples/` as starting points to create your own executable.

### Listening for Things

There are two ways to trigger actions on the bot:

- `Hear`: Trigger an action based on some criteria heard anywhere in the channel.
- `Respond`: Respond to a direct statement at the bot (e.g. "virbot show not shipped")

```go
// Create the bot (with adapter and name)
bot := victor.New("campfire", "ralph")

// Trigger anytime it hears the word "alot".
bot.Hear("alot", func(m victor.Message) {
    m.Room().Say("A LOT.")
})

// Trigger when someone talks directly to the bot in the following forms: 
//   - "ralph hi"
//   - "@ralph hi"
//   - "/hi"
bot.Respond("hi|hello|howdy", func(m victor.Message) {
    m.Reply(fmt.Sprintf("Hello, %s", m.User().Name()))
})

bot.Respond("bye", func(m victor.Message) {
    m.Reply(fmt.Sprintf("Goodbye, %s", m.User().Name()))
})

// Play a sound
bot.Respond("shipit", func(m victor.Message) {
    m.Room().Sound("nyan")
})

// Send a paste
bot.Respond("paste me", func(m victor.Message) {
    m.Room().Paste("/* Very important comment. */")
})

// Send a tweet URL
bot.Respond("gophers", func(m victor.Message) {
    m.Room().Tweet("https://twitter.com/brettbuddin/status/415852805623402496")
})

// Capture Parameters on trigger
bot.Respond("testsuite ([\\w-]+)/([\\w-]+):([\\w-]+)", func(m victor.Message) {
    params := m.Params()
    user   := params[0]
    repo   := params[1]
    branch := params[2]

    log.Printf("Running testsuite for user=%s repo=%s branch=%s\n", user, repo, branch)
    m.Reply("Running testsuite. Sit tight.")
})

bot.Run()
```
