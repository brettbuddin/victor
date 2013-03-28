package victor

import (
    "fmt"
    "strings"

    "github.com/brettbuddin/victor/shell"
    "github.com/edsrzf/fineline"
)

type Shell struct {
    brain *Brain
}

func NewShell(name string) *Shell {
    return &Shell{brain: NewBrain(name)}
}

func (s *Shell) Brain() *Brain {
    return s.brain
}

func (s *Shell) Run() {
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

        msg := &shell.Message{}
        msg.SetBody(command)

        ctx.SetMessage(msg)

        switch command {
        default:
            s.brain.Receive(ctx)
        case "close", "exit":
            fmt.Println("See ya!")
            return
        }

    }
}

func (s *Shell) Hear(expStr string, callback func(*Context)) {
    s.brain.Hear(expStr, callback)
}

func (s *Shell) Respond(expStr string, callback func(*Context)) {
    s.brain.Respond(expStr, callback)
}
