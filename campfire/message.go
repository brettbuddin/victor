package campfire

import (
    "encoding/json"
)

type Message struct {
    *Client

    id     int
    typ    string
    body   string
    roomId int
    userId int
}

type MessageWrapper struct {
    Message *Message `json:"message"`
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

func (m *Message) SetType(val string) {
    m.typ = val
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

func (m *Message) SetRoomId(val int) {
    m.roomId = val
}

func (m *Message) UserId() int {
    return m.userId
}

func (m *Message) SetUserId(val int) {
    m.userId = val
}

//
// JSON interface fulfillment
//

type messageData struct {
    Id     int    `json:"id,omitempty"`
    Type   string `json:"type"`
    Body   string `json:"body"`
    RoomId int    `json:"room_id,omitempty"`
    UserId int    `json:"user_id,omitempty"`
}

func (m *Message) MarshalJSON() ([]byte, error) {
    var data messageData

    data.Id = m.Id()
    data.Type = m.Type()
    data.Body = m.Body()

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

    m.SetId(actual.Id)
    m.SetType(actual.Type)
    m.SetBody(actual.Body)
    m.SetRoomId(actual.RoomId)
    m.SetUserId(actual.UserId)

    return nil
}
