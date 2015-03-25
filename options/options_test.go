package options

import (
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestStringer(t *testing.T) {
	options, _ := NewOptions("term", // signal
		[]string{}, // envs
		0,          // gracefulTimeout
		"json",     // configFormat
		"",         // configPath
		[]string{}, // pathsToTracks
		"",         // lockFile
		false,      // pty
		false,      // supervise
		false)      // restartOnConfigChanges

	assert.True(t, options.String() != "")
}

func TestIncorrectConfigFormat(t *testing.T) {
	_, err := NewOptions("term", // signal
		[]string{}, // envs
		0,          // gracefulTimeout
		"WTF",      // configFormat
		"",         // configPath
		[]string{}, // pathsToTracks
		"",         // lockFile
		false,      // pty
		false,      // supervise
		false)      // restartOnConfigChanges

	assert.NotNil(t, err)
}

func TestIncorrectSignalName(t *testing.T) {
	_, err := NewOptions("WTF", // signal
		[]string{}, // envs
		0,          // gracefulTimeout
		"json",     // configFormat
		"",         // configPath
		[]string{}, // pathsToTracks
		"",         // lockFile
		false,      // pty
		false,      // supervise
		false)      // restartOnConfigChanges

	assert.NotNil(t, err)
}
