package shell

type Message struct {
    id     int
    typ    string
    body   string
    roomId int
    userId int
}

func (m *Message) SetBody(val string) {
    m.body = val
}

func (m *Message) Id() int {
    return m.id
}

func (m *Message) Type() string {
    return m.typ
}

func (m *Message) Body() string {
    return m.body
}

func (m *Message) RoomId() int {
    return m.roomId
}

func (m *Message) UserId() int {
    return m.userId
}
