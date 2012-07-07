package victor

import (
    "fmt"
	"strings"

	"github.com/edsrzf/fineline"
)

type Shell struct {
    *Robot
}

func NewShell(robot *Robot) *Shell {
    return &Shell{Robot: robot}
}

func (self *Shell) Run() {
    reader := fineline.NewLineReader(nil)

    for {
        command, err := reader.Read()

		command = strings.TrimRight(command, "\n")

        if err != nil {
            break
        }

        msg  := &TextMessage{
            Body: command,

            Send: func(text string) {
                fmt.Println(text)    
            },
            Reply: func(text string) {
                fmt.Println("You: " + text)
            },
        }

        switch command {
            default:
                self.Receive(msg)
            case "close", "exit":
                self.Shutdown()
                return
        }

    }
}
