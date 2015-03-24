package lockfile

import (
	"fmt"
	"os"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

type Lock struct {
	name string
	file *os.File
}

func (l *Lock) String() string {
	return fmt.Sprintf("<Lock=(filename='%s', file='%v')>", l.name, l.file)
}

func (l *Lock) Acquire() (err error) {
	if l.file != nil {
		return fmt.Errorf("File %v is already acquired", l.file)
	}

	file, err := os.Create(l.name)
	if err != nil {
		log.WithField("filename", l.name).Info("Cannot create file.")
		return
	}

	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		log.WithField("lock", l).Info("Cannot acquire lock.")
		file.Close()
		return
	}

	l.file = file

	return
}

func (l *Lock) Release() error {
	defer func() {
		l.file.Close()
		os.Remove(l.name)
	}()

	err := syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
	if err != nil {
		log.WithFields(log.Fields{
			"filename":   l.name,
			"descriptor": int(l.file.Fd()),
		}).Error("Cannot release lock.")
	}
	l.file.Close()
	l.file = nil

	return err
}

func NewLock(filename string) *Lock {
	return &Lock{name: filename}
}
