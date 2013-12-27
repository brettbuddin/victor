package victor

import (
	"fmt"
	"github.com/brettbuddin/victor/util/google"
	"math/rand"
	"strconv"
	"time"
)

func defaults(robot *Robot) {
	robot.Respond("shutdown!", func(m Message) {
	    robot.Stop()
	})

	robot.Respond("ping", func(m Message) {
		m.Reply("pong!")
	})

	robot.Respond("roll( (\\d+))?", func(m Message) {
		defer recover()

		bound := 100
		val := m.Params()[1]

		if val != "" {
			var err error
			bound, err = strconv.Atoi(val)

			if err != nil {
				return
			}
		}

		rand.Seed(time.Now().UTC().UnixNano())
		random := rand.Intn(bound)
		m.Reply(fmt.Sprintf("rolled a %d of %d", random, bound))
	})

	robot.Respond("(image|img|gif|animate) (.*)", func(m Message) {
		gifOnly := (m.Params()[0] == "gif" || m.Params()[0] == "animate")

		result, err := google.ImageSearch(m.Params()[1], gifOnly)

		if err != nil {
			m.Room().Say("There was error making the request.")
			return
		}

		if result == "" {
			m.Room().Say("I didn't find anything.")
			return
		}

		m.Room().Say(result)
	})
}
