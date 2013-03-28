package campfire

import (
    "fmt"
    "encoding/json"
)

type User struct {
    *Client `json:"-"`

    id int
    typ string
    name string
    emailAddress string
    avatarUrl string
}

func (u *User) Show() (*User, error) {
    resp, err := u.Client.Get(fmt.Sprintf("/user/%d.json", u.Id()))

    if err != nil {
        return nil, err
    }

    var fetch struct {
        User *User `json:"user"`
    }

    err = json.Unmarshal(body(resp), &fetch)

    if err != nil {
        return nil, err
    }

    return fetch.User, nil
}

type userData struct {
    Id int      `json:"id"`
    Type string `json:"type"`
    Name string `json:"name"`
    EmailAddress string `json:"email_address"`
    AvatarUrl string `json:"avatar_url"`
}

func (u *User) MarshalJSON() ([]byte, error) {
    var data userData

    out, err := json.Marshal(data)

    if err != nil {
        return nil, err
    }

    return out, nil
}

func (u *User) UnmarshalJSON(data []byte) error {
    var actual userData

    err := json.Unmarshal(data, &actual)

    if err != nil {
        return err
    }

    u.id     = actual.Id
    u.typ    = actual.Type
    u.name   = actual.Name
    u.emailAddress = actual.EmailAddress
    u.avatarUrl    = actual.AvatarUrl

    return nil
}

func (u *User) Id() int {
    return u.id
}

func (u *User) Type() string {
    return u.typ
}

func (u *User) Name() string {
    return u.name
}

func (u *User) EmailAddress() string {
    return u.emailAddress
}

func (u *User) AvatarUrl() string {
    return u.avatarUrl
}
