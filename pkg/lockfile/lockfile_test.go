package lockfile

import "testing"

func TestLocking(t *testing.T) {
    file, _ := NewTemp("/tmp")
    done := make(chan bool)

    for i := 0; i < 10; i++ {
        go func() {
            for j := 0; j < 1000; j++ {
                file.Lock()
                file.Unlock()
            }

            done <- true
        }()
    }

    for i := 0; i < 10; i++ {
        <-done
    }
}
