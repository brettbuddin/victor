package ssh

import (
    "errors"
    "io"
    "crypto"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
)

type Keychain struct {
    keys []interface{}
}

func (k *Keychain) Add(privateKey []byte) error {
    block, _ := pem.Decode(privateKey)

    if block == nil {
        return errors.New("no PEM data found")
    }

    key, err := x509.ParsePKCS1PrivateKey(block.Bytes)

    if err != nil {
        return err
    }

    k.keys = append(k.keys, key)

    return nil
}

func (k *Keychain) Key(i int) (interface{}, error) {
    if i < 0 || i >= len(k.keys) {
        return nil, nil
    }
    
    switch key := k.keys[i].(type) {
    case *rsa.PrivateKey:
        return &key.PublicKey, nil
    }
    panic("unknown key type")
}

func (k *Keychain) Sign(i int, rand io.Reader, data []byte) ([]byte, error) {
    h := crypto.SHA1.New()
    h.Write(data)
    digest := h.Sum(nil)

    switch key := k.keys[i].(type) {
    case *rsa.PrivateKey:
        return rsa.SignPKCS1v15(rand, key, crypto.SHA1, digest)
    }

    return nil, errors.New("unknown key type")
}
