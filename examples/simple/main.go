package main

import (
	"fmt"
	"github.com/brettbuddin/victor"
	"os"
	"os/signal"
)

func main() {
	bot := victor.New("shell", "victor", ":8000")

	bot.HandleFunc(bot.Direct("hello|hi|howdy"), func(s victor.State) {
		s.Chat().Send(s.Message().ChannelID(), fmt.Sprintf("Hello, %s", s.Message().UserName()))
	})

	go bot.Run()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	<-sigs

	bot.Stop()
}
