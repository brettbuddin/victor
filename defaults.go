package victor

import (
    "github.com/brettbuddin/victor/utils/google"
    "strconv"
    "fmt"
    "time"
    "math/rand"
)

func registerDefaultAbilities(brain *Brain) {
    brain.Respond("campfire id", func(ctx *Context) {
        id := strconv.Itoa(ctx.Message().UserId())
        ctx.Reply(id)
    })

    brain.Respond("ping", func(ctx *Context) {
        ctx.Reply("pong!")
    })

    brain.Respond("roll( (\\d+))?", func(ctx *Context) {
        defer recover()

        bound      := 100
        val        := ctx.Matches()[2]

        if val != "" {
            var err error
            bound, err = strconv.Atoi(val)

            if err != nil {
                return
            }
        }

        rand.Seed(time.Now().UTC().UnixNano())
        random := rand.Intn(bound)
        ctx.Reply(fmt.Sprintf("rolled a %d of %d", random, bound))
    })

    brain.Respond("(image|img|gif|animate) (.*)", func(ctx *Context) {
        gifOnly := (ctx.Matches()[1] == "gif" || ctx.Matches()[1] == "animate")

        result, err := google.ImageSearch(ctx.Matches()[2], gifOnly)

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
