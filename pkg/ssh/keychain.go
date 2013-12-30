package ssh

import (
	"code.google.com/p/go.crypto/ssh"
	"io"
)

type Keychain struct {
	keys []ssh.Signer
}

func (k *Keychain) Add(privateKey []byte) error {
	key, err := ssh.ParsePrivateKey(privateKey)

	if err != nil {
		return err
	}

	k.keys = append(k.keys, key)

	return nil
}

func (k *Keychain) Key(i int) (ssh.PublicKey, error) {
	if i < 0 || i >= len(k.keys) {
		return nil, nil
	}

	return k.keys[i].PublicKey(), nil
}

func (k *Keychain) Sign(i int, rand io.Reader, data []byte) ([]byte, error) {
	return k.keys[i].Sign(rand, data)
}
