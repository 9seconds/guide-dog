package lockfile

import (
	"io/ioutil"
	"os"
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func makeTempFile() string {
	tempFile, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}

	return tempFile.Name()
}

func TestAcquire(t *testing.T) {
	fileName := makeTempFile()
	defer os.Remove(fileName)

	lock := NewLock(fileName)
	defer lock.Release()

	assert.Nil(t, lock.Acquire())
	assert.NotNil(t, lock.Acquire())
}

func TestReleaseOK(t *testing.T) {
	fileName := makeTempFile()
	defer os.Remove(fileName)

	lock := NewLock(fileName)
	defer lock.finish()

	assert.Nil(t, lock.Acquire())
	assert.Nil(t, lock.Release())

	_, err := os.Stat(fileName)
	assert.Nil(t, err)
}

func TestReleaseAbsentLock(t *testing.T) {
	lock := NewLock("WTF")

	assert.NotNil(t, lock.Release())
}
