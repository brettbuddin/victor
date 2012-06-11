package campfire

import (
    "net/url"
    "net/http"
    "bytes"
    "io/ioutil"
    "encoding/json"
    "log"
)

type Client struct {
    stream  *Stream
    account string
    token   string
}

func NewClient(account, token string) *Client {
    return &Client{
        account: account,
        token: token,
    }
}

func (self *Client) Room(id int) *Room {
    return &Room{client: self, Id: id}
}

func (self *Client) Me() (*User, error) {
    resp, err := self.Get("/users/me")

    if err != nil {
        return nil, err
    }

    out, err := ioutil.ReadAll(resp.Body)

    if err != nil {
        return nil, err
    }

    var fetch struct {
        User *User `json:"user"`
    }

    err = json.Unmarshal(out, &fetch)

    if err != nil {
        return nil, err
    }

    return fetch.User, nil
}

func (self *Client) Get(path string) (*http.Response, error) {
    return self.request("GET", path, "")
}

func (self *Client) Post(path string, body string) (*http.Response, error) {
    return self.request("POST", path, body)
}

func (self *Client) Stream(roomId int, channel chan *Message) {
    self.stream = NewStream(self.token, self.Room(roomId))
    self.stream.Connect()

    resp, _ := self.stream.Connect()

    if resp.StatusCode != 200 {
        log.Fatal(resp.Status)
    }

    go self.stream.Read(resp, channel)
}

func (self *Client) request(method string, path string, body string) (*http.Response, error) {
    url        := new(url.URL) 
    url.Scheme = "https"
    url.Host   = self.account + ".campfirenow.com"
    url.Path   = path

    client    := &http.Client{}

    req, err := http.NewRequest(method, url.String(), nil)
    req.Header.Add("Content-Type", "application/json")
    req.SetBasicAuth(self.token, "X")

    if method == "POST" {
        req.Body = closer{bytes.NewBufferString(body)}
        req.ContentLength = int64(len(body))
    }

    if err != nil {
        return nil, err
    }
    
    return client.Do(req)
}
