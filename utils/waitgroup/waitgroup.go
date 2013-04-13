package waitgroup

import (
    "sync"
)

type Wrapper struct {
    sync.WaitGroup
}

func (w *Wrapper) Wrap(fn func()) {
    w.Add(1)
    go func() {
        fn()
        w.Done()
    }()
}
