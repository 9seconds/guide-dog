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
		false,      // restartOnConfigChanges
		[]string{}) // exitCodes

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
		false,      // restartOnConfigChanges
		[]string{}) // exitCodes

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
		false,      // restartOnConfigChanges
		[]string{}) // exitCodes

	assert.NotNil(t, err)
}

func TestIncorrectExitCodes(t *testing.T) {
	_, err := NewOptions("term", // signal
		[]string{},      // envs
		0,               // gracefulTimeout
		"json",          // configFormat
		"",              // configPath
		[]string{},      // pathsToTracks
		"",              // lockFile
		false,           // pty
		false,           // supervise
		false,           // restartOnConfigChanges
		[]string{"ggg"}) // exitCodes

	assert.NotNil(t, err)
}

func TestCorrectExitCodes(t *testing.T) {
	options, err := NewOptions("term", // signal
		[]string{}, // envs
		0,          // gracefulTimeout
		"json",     // configFormat
		"",         // configPath
		[]string{}, // pathsToTracks
		"",         // lockFile
		false,      // pty
		false,      // supervise
		false,      // restartOnConfigChanges
		[]string{"1", "2", "1"}) // exitCodes

	assert.Nil(t, err)
	assert.Equal(t, len(options.ExitCodes), 2)
	assert.True(t, options.ExitCodes[1])
	assert.True(t, options.ExitCodes[2])
}
