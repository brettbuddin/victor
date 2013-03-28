package campfire

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net"
    "net/http"
    "net/http/httputil"
    "net/url"
    "crypto/tls"
)

type Client struct {
    host    string
    token   string
}

func NewClient(subdomain, token string) *Client {
    return &Client{
        host:   subdomain + ".campfirenow.com",
        token:  token,
    }
}

func (c *Client) Room(id int) *Room {
    return &Room{Client: c, id: id}
}

func (c *Client) User(id int) *User {
    return &User{Client: c, id: id}
}

func (c *Client) Account(id int) *Account {
    return &Account{Client: c, Id: id}
}

func (c *Client) Me() (*User, error) {
    resp, err := c.Get("/users/me")

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

func (c *Client) Get(path string) (*http.Response, error) {
    return c.request("GET", path, []byte(""))
}

func (c *Client) Post(path string, body []byte) (*http.Response, error) {
    return c.request("POST", path, body)
}

func (c *Client) request(method, path string, body []byte) (*http.Response, error) {
    url := &url.URL{
        Scheme: "https",
        Host:   c.host,
        Path:   path,
    }

    conn, err := net.Dial("tcp", url.Host + ":443")

    if err != nil {
        return nil, err
    }

    ssl        := tls.Client(conn, nil)
    httpClient := httputil.NewClientConn(ssl, nil)

    req, err := http.NewRequest(method, url.String(), bytes.NewBuffer(body))
    req.Header.Add("Content-Type", "application/json")
    req.SetBasicAuth(c.token, "X")

    if method == "POST" {
        req.ContentLength = int64(len(body))
    }

    if err != nil {
        return nil, err
    }

    return httpClient.Do(req)
}

func body(resp *http.Response) []byte {
    out, err := ioutil.ReadAll(resp.Body)

    if err != nil {
        return []byte{}
    }

    return out
}
