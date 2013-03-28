package campfire

import (
    "encoding/json"
    "errors"
    "fmt"
)

type Room struct {
    *Client

    id    int
    name  string
    users []*User
}

func (r *Room) Show() (*Room, error) {
    resp, err := r.Client.Get(fmt.Sprintf("/room/%d.json", r.Id()))

    if err != nil {
        return nil, err
    }

    var fetch struct {
        Room *Room `json:"room"`
    }

    err = json.Unmarshal(body(resp), &fetch)

    if err != nil {
        return nil, err
    }

    return fetch.Room, nil
}

func (r *Room) Join() error {
    resp, err := r.Client.Post(fmt.Sprintf("/room/%d/join", r.Id()), []byte(""))

    if err == nil && resp.StatusCode == 200 {
        return nil
    }

    return errors.New("Could not join room.")
}

func (r *Room) Say(message string) {
    r.sendMessage(message, "TextMessage")
}

func (r *Room) Paste(message string) {
    r.sendMessage(message, "PasteMessage")
}

func (r *Room) Sound(name string) {
    r.sendMessage(name, "SoundMessage")
}

func (r *Room) Tweet(message string) {
    r.sendMessage(message, "TweetMessage")
}

func (r *Room) sendMessage(message, typ string) {
    msg := &MessageWrapper{Message: &Message{typ: typ, body: message}}
    buf, _ := json.Marshal(msg)

    r.Client.Post(fmt.Sprintf("/room/%d/speak", r.Id()), buf)
}

func (r *Room) Stream(incoming chan *Message) {
    stream := NewStream(r.Client, r)
    stream.Start(incoming)
}

// INTERFACE FULFILLMENT

type roomData struct {
    Id    int     `json:"id"`
    Name  string  `json:"name"`
    Users []*User `json:"users"`
}

func (r *Room) MarshalJSON() ([]byte, error) {
    var data roomData

    out, err := json.Marshal(data)

    if err != nil {
        return nil, err
    }

    return out, nil
}

func (r *Room) UnmarshalJSON(data []byte) error {
    var actual roomData

    err := json.Unmarshal(data, &actual)

    if err != nil {
        return err
    }

    r.id     = actual.Id
    r.name   = actual.Name
    r.users  = actual.Users

    return nil
}

func (r *Room) Id() int {
    return r.id
}

func (r *Room) Name() string {
    return r.name
}

func (r *Room) Users() []*User {
    return r.users
}

