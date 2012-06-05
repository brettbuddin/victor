package victor

type TextMessage struct {
    Id        int
    Body      string
    CreatedAt string

    Send  func(text string)
    Reply func(text string)
    Topic func(text string)
}
