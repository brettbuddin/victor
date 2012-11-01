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

        ctx := &Context{
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

        ctx.SetMessage(&Message{Body: command})

        switch command {
        default:
            self.brain.Receive(ctx)
        case "close", "exit":
            fmt.Println("See ya!")
            return
        }

    }
}

func (self *Shell) Hear(expStr string, callback func(*Context)) {
    self.brain.Hear(expStr, callback)
}

func (self *Shell) Respond(expStr string, callback func(*Context)) {
    self.brain.Respond(expStr, callback)
}
