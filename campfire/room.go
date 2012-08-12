package campfire

import (
    "encoding/json"
    "errors"
    "io/ioutil"
    "strconv"
)

type Room struct {
    client *Client

    Id    int     `json:"id"`
    Name  string  `json:"name"`
    Users []*User `json:"users"`
}

func (self *Room) Show() (*Room, error) {
    resp, err := self.client.Get("/room/" + strconv.Itoa(self.Id) + ".json")

    if err != nil {
        return nil, err
    }

    out, err := ioutil.ReadAll(resp.Body)

    if err != nil {
        return nil, err
    }

    var fetch struct {
        Room *Room `json:"room"`
    }

    err = json.Unmarshal(out, &fetch)

    if err != nil {
        return nil, err
    }

    return fetch.Room, nil
}

func (self *Room) Join() error {
    resp, err := self.client.Post("/room/"+strconv.Itoa(self.Id)+"/join", []byte(""))

    if err == nil && resp.StatusCode == 200 {
        return nil
    }

    return errors.New("Could not join room.")
}

func (self *Room) Stream(channel chan *Message) {
    self.client.Stream(self.Id, channel)
}

func (self *Room) Say(message string) {
    msg := &MessageWrapper{Message: &Message{Type: "TextMessage", Body: message}}
    buf, _ := json.Marshal(msg)

    self.client.Post("/room/"+strconv.Itoa(self.Id)+"/speak", buf)
}

func (self *Room) Paste(message string) {
    msg := &MessageWrapper{Message: &Message{Type: "PasteMessage", Body: message}}
    buf, _ := json.Marshal(msg)

    self.client.Post("/room/"+strconv.Itoa(self.Id)+"/speak", buf)
}

func (self *Room) Sound(name string) {
    msg := &MessageWrapper{Message: &Message{Type: "SoundMessage", Body: name}}
    buf, _ := json.Marshal(msg)

    self.client.Post("/room/"+strconv.Itoa(self.Id)+"/speak", buf)
}
