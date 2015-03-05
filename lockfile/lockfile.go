package lockfile

import (
	"fmt"
	"os"
	"syscall"
)

type Lock struct {
	name string
	file *os.File
}

func (l *Lock) String() string {
	return fmt.Sprintf("<Lock=(filename='%s', file='%v')>", l.name, l.file)
}

func (l *Lock) Acquire() (err error) {
	if l.file == nil {
		return
	}

	file, err := os.Create(l.name)
	if err != nil {
		return
	}

	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		file.Close()
		return
	}

	l.file = file

	return
}

func (l *Lock) Release() error {
	defer l.file.Close()

	err := syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
	l.file.Close()
	l.file = nil

	return err
}

func NewLock(filename string) *Lock {
	return &Lock{name: filename}
}
