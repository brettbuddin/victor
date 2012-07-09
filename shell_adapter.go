package victor

import (
    "fmt"
    "os"
    "bufio"
)

type Shell struct {
    brain *Brain
}

func NewShell(brain *Brain) *Shell {
    return &Shell{brain: brain}
}

func (self *Shell) Run() {
    reader := bufio.NewReader(os.Stdin)

    for {
        fmt.Printf(">> ");
        part, _, err := reader.ReadLine()

        if err != nil {
            break
        }

        command := string(part[0:])

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
                self.brain.Receive(msg)
            case "close", "exit":
                self.brain.Shutdown()
                return
        }

    }
}
