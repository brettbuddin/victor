package shell

import (
	"bufio"
	"fmt"
	"github.com/brettbuddin/victor/pkg/chat"
	"os"
)

func init() {
	chat.Register("shell", func(r chat.Robot) chat.Adapter {
		return &Adapter{r, make(chan bool)}
	})
}

type Adapter struct {
	robot chat.Robot
	stop  chan bool
}

func (a *Adapter) Run() {
	reader := bufio.NewReader(os.Stdin)

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
			return
		case line := <-lines:
			a.robot.Receive(&message{string(line)})
		}
	}
}

func (a *Adapter) Send(channelID, msg string) {
	fmt.Println("SEND:", msg)
}

func (a *Adapter) Stop() {
	a.stop <- true
	close(a.stop)
}

type message struct {
	text string
}

func (m *message) UserID() string {
	return "0"
}

func (m *message) UserName() string {
	return "Meathead"
}

func (m *message) ChannelID() string {
	return "0"
}

func (m *message) ChannelName() string {
	return "#ops"
}

func (m *message) Text() string {
	return m.text
}
