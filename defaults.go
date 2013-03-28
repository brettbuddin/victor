package victor

import (
    "github.com/brettbuddin/victor/utils/google"
    "strconv"
)

func registerDefaultAbilities(brain *Brain) {
    brain.Respond("campfire id", func(ctx *Context) {
        id := strconv.Itoa(ctx.Message().Id())
        ctx.Reply(id)
    })

    brain.Respond("ping", func(ctx *Context) {
        ctx.Reply("pong!")
    })

    brain.Respond("(image|img|gif|animate) (.*)", func(ctx *Context) {
        gifOnly := (ctx.Matches()[0] == "gif" || ctx.Matches()[0] == "animate")

        result, err := google.ImageSearch(ctx.Matches()[1], gifOnly)

        if err != nil {
            ctx.Send("There was error making the request.")
            return
        }

        if result == "" {
            ctx.Send("I didn't find anything.")
            return
        }

        ctx.Send(result)
    })

    brain.Respond("users (list|show)", func(ctx *Context) {
        list := ""

        for _, user := range brain.Users() {
            list += user.Name() + "\n"
        }

        ctx.Paste(list)
    })
}
