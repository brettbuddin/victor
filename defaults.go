package victor

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func defaults(robot Robot) {
	robot.HandleFunc(robot.Direct("ping"), func(s State) {
		s.Chat().Send(s.Message().ChannelID(), "pong!")
	})

	robot.HandleFunc(robot.Direct("roll(\\s(\\d+))?"), func(s State) {
		defer recover()

		bound := 100
		val := s.Params()[1]

		if val != "" {
			var err error
			bound, err = strconv.Atoi(val)

			if err != nil {
				return
			}
		}

		rand.Seed(time.Now().UTC().UnixNano())
		random := rand.Intn(bound)

		msg := fmt.Sprintf("%s rolled a %d of %d", s.Message().UserName(), random, bound)
		s.Chat().Send(s.Message().ChannelID(), msg)
	})
}
