package campfire

import (
    "bufio"
    "crypto/tls"
    "encoding/json"
    "log"
    "net"
    "net/http"
    "net/http/httputil"
    "net/url"
    "strconv"
    "time"
)

type Stream struct {
    connection *httputil.ClientConn
    token      string
    room       *Room
}

func NewStream(token string, room *Room) *Stream {
    return &Stream{token: token, room: room}
}

func (self *Stream) Connect() (*http.Response, error) {
    url := new(url.URL)
    url.Scheme = "https"
    url.Host = "streaming.campfirenow.com"
    url.Path = "/room/" + strconv.Itoa(self.room.Id) + "/live.json"

    conn, err := net.Dial("tcp", url.Host+":443")

    if err != nil {
        return nil, err
    }

    ssl := tls.Client(conn, nil)
    client := httputil.NewClientConn(ssl, nil)

    req, err := http.NewRequest("GET", url.String(), nil)

    req.SetBasicAuth(self.token, "X")

    resp, err := client.Do(req)

    if err != nil {
        log.Print("Couldn't initiate stream.")
        return nil, err
    }

    return resp, nil
}

func (self *Stream) Read(resp *http.Response, channel chan *Message) {
    reader := bufio.NewReader(resp.Body)

    for {
        line, err := reader.ReadBytes('\r')

        if err != nil {
            time.Sleep(6)

            log.Print("Reconnecting...")
            self.room.Join()
            resp, err = self.Connect()

            if err != nil || resp.StatusCode != 200 {
                panic("Could not reconnect.")
                continue
            }

            reader = bufio.NewReader(resp.Body)
            continue
        }

        var msg Message

        err = json.Unmarshal(line, &msg)

        if err != nil {
            continue
        }

        channel <- &msg
    }
}
