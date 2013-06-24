package main

import (
    "github.com/brettbuddin/victor"
    "fmt"
)

func main() {
    bot, err := victor.New("shell", "bot")

    if err != nil {
        fmt.Println(err)
    }

    bot.Respond("hello|hi|howdy", func(m victor.Message) {
        m.Room().Say(fmt.Sprintf("Hello, %s", m.User().Name()))
    })

    bot.Run()
}
