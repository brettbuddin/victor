package campfire

import (
    "fmt"
    "encoding/json"
)

type User struct {
    *Client

    id int
    typ string
    name string
    emailAddress string
    avatarUrl string
}

func (u *User) Id() int {
    return u.id
}

func (u *User) SetId(val int) {
    u.id = val
}

func (u *User) Type() string {
    return u.typ
}

func (u *User) SetType(val string) {
    u.typ = val
}

func (u *User) Name() string {
    return u.name
}

func (u *User) SetName(val string) {
    u.name = val
}

func (u *User) EmailAddress() string {
    return u.emailAddress
}

func (u *User) SetEmailAddress(val string) {
    u.emailAddress = val
}

func (u *User) AvatarUrl() string {
    return u.avatarUrl
}

func (u *User) SetAvatarUrl(val string) {
    u.avatarUrl = val
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

//
// JSON interface fulfillment
//

type userData struct {
    Id int      `json:"id"`
    Type string `json:"type"`
    Name string `json:"name"`
    EmailAddress string `json:"email_address"`
    AvatarUrl string `json:"avatar_url"`
}

func (u *User) MarshalJSON() ([]byte, error) {
    var data userData

    data.Id     = u.Id()
    data.Type   = u.Type()
    data.Name   = u.Name()
    data.EmailAddress = u.EmailAddress()
    data.AvatarUrl    = u.AvatarUrl()

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

    u.SetId(actual.Id)
    u.SetType(actual.Type)
    u.SetName(actual.Name)
    u.SetEmailAddress(actual.EmailAddress)
    u.SetAvatarUrl(actual.AvatarUrl)

    return nil
}
