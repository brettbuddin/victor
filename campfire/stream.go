package campfire

import (
    "bufio"
    "encoding/json"
    "net/http"
    "fmt"
    "time"
)

type Stream struct {
    parentClient *Client
    room         *Room
    exit         bool
}

func NewStream(parent *Client, room *Room) *Stream {
    return &Stream{parentClient: parent, room: room}
}

func (s *Stream) Start(incoming chan *Message) {
    resp, err := s.connect()

    if err != nil || resp.StatusCode != 200 {
        return
    }

    go s.consume(resp, incoming)
}

func (s *Stream) Stop() {
    s.exit = true
}

func (s *Stream) connect() (*http.Response, error) {
    client := NewClient("streaming", s.parentClient.token)
    return client.Get(fmt.Sprintf("/room/%d/live.json", s.room.Id()))
}

func (s *Stream) consume(resp *http.Response, channel chan *Message) {
    reader := bufio.NewReader(resp.Body)

    for {
        line, err := reader.ReadBytes('\r')

        if s.exit {
            return
        }

        if err != nil {
            time.Sleep(6)

            s.room.Join()
            resp, err = s.connect()

            if err != nil || resp.StatusCode != 200 {
                panic("Could not reconnect.")
                continue
            }

            reader = bufio.NewReader(resp.Body)
            continue
        }

        var msg *Message

        err = json.Unmarshal(line, msg)

        if err != nil {
            continue
        }

        msg.Client = s.parentClient

        channel <- msg
    }
}

