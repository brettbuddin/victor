package shell

import (
	"bufio"
	"fmt"
	"github.com/brettbuddin/victor/adapter"
	"os"
)

func init() {
	adapter.Register("shell", func(b adapter.Brain) adapter.Adapter {
		return &Shell{Brain: b}
	})
}

type Shell struct {
	adapter.Brain
}

func (s *Shell) Listen(messages chan adapter.Message) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Type your commands (type \"exit\" to exit):")

	for {
		line, _, err := reader.ReadLine()

		if err != nil {
			continue
		}

		switch string(line) {
		case "exit":
			return nil
		default:
			messages <- &Message{
				body: string(line),
				user: &User{0, "You"},
				room: &Room{0, "Chat City"},
			}
		}
	}

	return nil
}
