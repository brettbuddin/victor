package shell

import (
	"bufio"
	"fmt"
	"github.com/brettbuddin/victor/adapter"
	"os"
)

func init() {
	adapter.Register("shell", func(adapter.Brain) adapter.Adapter {
		return &Adapter{
		    stop: make(chan bool),
		}
	})
}

type Adapter struct {
    stop chan bool
}

func (a *Adapter) Listen(messages chan adapter.Message) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Type your commands (type \"exit\" to exit):")

    lines := make(chan string)
	go func() {
        for {
            if line, _, err := reader.ReadLine(); err == nil {
                lines <- string(line)
            }
        }
    }()

	for {
		select {
		case <-a.stop:
		    close(messages)
			return nil
		case line := <-lines:
			messages <- &Message{
				body:   string(line),
				params: []string{},
				user:   User{0, "You"},
				room:   Room{0, "Chat City"},
			}
		}
	}
}

func (a *Adapter) Stop() {
    a.stop <- true
    close(a.stop)
}
