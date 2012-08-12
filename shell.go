package victor

import (
    "fmt"
    "strings"

    "github.com/edsrzf/fineline"
)

type Shell struct {
    brain *Brain
}

func NewShell(name string) *Shell {
    return &Shell{brain: NewBrain(name)}
}

func (self *Shell) Run() {
    reader := fineline.NewLineReader(nil)

    for {
        command, err := reader.Read()

        command = strings.TrimRight(command, "\n")

        if err != nil {
            break
        }

        msg := &TextMessage{
            Body: command,

            Send: func(text string) {
                fmt.Println(text)
            },
            Reply: func(text string) {
                fmt.Println("You: " + text)
            },
            Paste: func(text string) {
                fmt.Println("You: " + text)
            },
            Sound: func(name string) {
                fmt.Println("Plays sound: " + name)
            },
        }

        switch command {
        default:
            self.brain.Receive(msg)
        case "close", "exit":
            self.brain.Shutdown()
            return
        }

    }
}

func (self *Shell) Hear(expStr string, callback func(*TextMessage)) {
    self.brain.Hear(expStr, callback)
}

func (self *Shell) Respond(expStr string, callback func(*TextMessage)) {
    self.brain.Respond(expStr, callback)
}
