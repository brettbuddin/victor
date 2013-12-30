package lockfile

import (
	"os"
	"syscall"
	"io/ioutil"
)

const (
	LOCK_FLAGS   = syscall.LOCK_EX | syscall.LOCK_NB
	UNLOCK_FLAGS = syscall.LOCK_UN
)

type LockFile struct {
	*os.File
}

func New(filepath string) (*LockFile, error) {
	file, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}

	return &LockFile{file}, nil
}

func NewTemp(dirPath string) (*LockFile, error) {
    file, err := ioutil.TempFile(dirPath, "")
	if err != nil {
		return nil, err
	}

	return &LockFile{file}, nil
}

func (l *LockFile) Lock() error {
	if err := syscall.Flock(int(l.File.Fd()), LOCK_FLAGS); err != nil {
		l.File.Close()
		return err
	}

	return nil
}

func (l *LockFile) Unlock() error {
	defer l.File.Close()
	if err := syscall.Flock(int(l.File.Fd()), UNLOCK_FLAGS); err != nil {
		return err
	}

	return nil
}

func (l *LockFile) Do(f func()) {
	l.Lock()
	defer l.Unlock()
	f()
}
