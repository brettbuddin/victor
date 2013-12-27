package main

import (
	"fmt"
	"github.com/brettbuddin/victor"
	"os"
	"os/signal"
)

func main() {
	bot, err := victor.New("shell", "bot")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(bot.Respond("(?P<phrase>hello|hi|howdy)", func(m victor.Message) {
		fmt.Println(m.Params())
		//m.Room().Say(fmt.Sprintf("Hello, %s", m.User().Name()))
	}))

	signals(bot).Run()
}

func signals(bot *victor.Robot) *victor.Robot {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	go func() {
		<-sigs
		bot.Stop()
	}()

	return bot
}
