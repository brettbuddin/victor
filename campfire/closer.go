package campfire

import "io"

type closer struct {
    io.Reader
}

func (closer) Close() error {
    return nil
}
