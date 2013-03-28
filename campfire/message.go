package campfire

import (
    "encoding/json"
)

type Message struct {
    *Client `json:"-"`

    id int
    typ string
    body string
    roomId int
    userId int
}

// INTERFACE FULFILLMENT

type messageData struct {
    Id int      `json:"id"`
    Type string `json:"type"`
    Body string `json:"body"`
    RoomId  int `json:"room_id"`
    UserId  int `json:"user_id"`
}

func (m *Message) MarshalJSON() ([]byte, error) {
    var data messageData

    out, err := json.Marshal(data)

    if err != nil {
        return nil, err
    }

    return out, nil
}

func (m *Message) UnmarshalJSON(data []byte) error {
    var actual messageData

    err := json.Unmarshal(data, &actual)

    if err != nil {
        return err
    }

    m.id     = actual.Id
    m.typ    = actual.Type
    m.body   = actual.Body
    m.roomId = actual.RoomId
    m.userId = actual.UserId

    return nil
}

type MessageWrapper struct {
    Message *Message
}

func (m *Message) Id() int {
    return m.id
}

func (m *Message) SetId(val int) {
    m.id = val
}

func (m *Message) Type() string {
    return m.typ
}

func (m *Message) Body() string {
    return m.body
}

func (m *Message) SetBody(val string) {
    m.body = val
}

func (m *Message) RoomId() int {
    return m.roomId
}

func (m *Message) UserId() int {
    return m.userId
}
