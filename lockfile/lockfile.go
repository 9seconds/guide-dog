// Package lockfile defines a flock(2)-based implementation of the file lock.
package lockfile

import (
	"fmt"
	"os"
	"sync"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

// File* constants family defines flags which are used to open lock files.
const (
	FileOpenFlags   = os.O_WRONLY
	FileCreateFlags = os.O_CREATE | os.O_EXCL | FileOpenFlags
)

// FileMode defines default permission for file lock open.
const FileMode = os.FileMode(0666)

// Lock file is a thin wrapper around os.File to give user a possibility
// to use file as a lock (using flock(2))
type Lock struct {
	name           string
	fileWasCreated bool
	openLock       *sync.Mutex
	file           *os.File
}

func (l *Lock) String() string {
	return fmt.Sprintf("<Lock=(filename='%s', fileWasCreated=%t, file='%v')>",
		l.name,
		l.fileWasCreated,
		l.file)
}

// Acquire file lock. Returns error if acquiring failed, nil otherwise.
func (l *Lock) Acquire() (err error) {
	if l.file != nil {
		return fmt.Errorf("File %v is already acquired", l.file)
	}

	if err = l.open(); err != nil {
		log.WithField("filename", l.name).Info("Cannot open file.")
		l.finish()
		return
	}

	err = syscall.Flock(int(l.file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		log.WithField("lock", l).Info("Cannot acquire lock.")
		l.finish()
		return
	}

	return
}

// Release file lock. Returns error if something went wrong, nil otherwise.
func (l *Lock) Release() error {
	defer l.finish()

	err := syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
	if err != nil {
		log.WithFields(log.Fields{
			"filename":   l.name,
			"descriptor": int(l.file.Fd()),
		}).Error("Cannot release lock.")
	}

	return err
}

// open correctly opens file with different modes.
func (l *Lock) open() (err error) {
	l.openLock.Lock()
	defer l.openLock.Unlock()

	l.fileWasCreated = false

	flags := FileOpenFlags
	if _, err := os.Stat(l.name); os.IsNotExist(err) {
		flags = FileCreateFlags
	}

	file, err := os.OpenFile(l.name, flags, FileMode)
	if err != nil && flags == FileCreateFlags {
		flags = FileOpenFlags
		file, err = os.OpenFile(l.name, flags, FileMode)
	}

	if err == nil {
		l.file = file
		l.fileWasCreated = flags == FileCreateFlags
	}

	return
}

// finish just cleans up lock file.
func (l *Lock) finish() {
	if l.file != nil {
		l.file.Close()
		if l.fileWasCreated {
			os.Remove(l.name)
		}
	}
}

// NewLock returns new lock instance. Argument is a path to the lock file.
// Please remember that file would be truncated so it should be file path
// for absent file or something which has insensitive content.
func NewLock(filename string) *Lock {
	return &Lock{name: filename, fileWasCreated: false, openLock: new(sync.Mutex)}
}
