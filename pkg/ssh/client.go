package ssh

import (
	"bytes"
	"code.google.com/p/go.crypto/ssh"
)

type Client struct {
	Config *ssh.ClientConfig
	Client *ssh.ClientConn
	Addr   string
}

func NewClient(user string, addr string, keys *Keychain) *Client {
	return &Client{
		Addr: addr,
		Config: &ssh.ClientConfig{
			User: user,
			Auth: []ssh.ClientAuth{
				ssh.ClientAuthKeyring(keys),
			},
		},
	}
}

func (c *Client) Connect() (err error) {
	var client *ssh.ClientConn

	client, err = ssh.Dial("tcp", c.Addr, c.Config)
	if err != nil {
		return
	}

	c.Client = client

	return
}

func (c *Client) Run(cmd string) (out bytes.Buffer, err error) {
	session, err := c.Client.NewSession()
	if err != nil {
		return
	}
	defer session.Close()

	session.Stdout = &out

	if err = session.Run(cmd); err != nil {
		return
	}

	return
}
