package victor

import (
    "github.com/brettbuddin/victor/utils/google"
)

func registerDefaultAbilities(brain *Brain) {
    brain.Respond("ping", func(msg *TextMessage) {
        msg.Reply("pong!")
    })

    brain.Respond("(image|img) (.*)", func(msg *TextMessage) {
        result, err := google.ImageSearch(msg.Matches()[3])

        if err != nil {
            return
        }

        if result == "" {
            msg.Send("I didn't find anything.")
            return
        }

        msg.Send(result)
    })

    brain.Respond("(list|show) users", func(msg *TextMessage) {
        list := ""

        for _, user := range brain.Users() {
            list += user.Name + "\n"
        }

        msg.Paste(list)
    })
}
