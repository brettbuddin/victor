package campfire

import (
    "encoding/json"
    "strconv"
    "io/ioutil"
    "errors"
)

type Room struct {
    Client *Client

    Id int `json:"id"`
    Name string `json:"name"`
    Users []*User `json:"users"`
}

func (self *Room) Show() (*Room, error) {
    resp, err := self.Client.Get("/room/" + strconv.Itoa(self.Id) + ".json")

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
    resp, err := self.Client.Post("/room/" + strconv.Itoa(self.Id) + "/join", "")

    if err == nil && resp.StatusCode == 200 {
        return nil
    }

    return errors.New("Could not join room.")
}

func (self *Room) Stream(channel chan *Message) {
    self.Client.Stream(self.Id, channel)
}

func (self *Room) Say(message string) {
    msg    := &MessageWrapper{Message: &Message{Type: "TextMessage", Body: message}}
    buf, _ := json.Marshal(msg)

    self.Client.Post("/room/" + strconv.Itoa(self.Id) + "/speak", string(buf))
}
