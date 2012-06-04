package campfire

type Message struct {
    Id         int `json:"id"`
    RoomId     int `json:"room_id"`
    UserId     int `json:"user_id"`

    Type       string `json:"type"`
    CreatedAt  string `json:"created_at"`
    Body       string `json:"body"`
}

type MessageWrapper struct {
    Message *Message `json:"message"`
}
