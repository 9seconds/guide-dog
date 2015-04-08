package environment

import (
	"io/ioutil"
	"os"
	"testing"

	assert "github.com/stretchr/testify/assert"

	opts "github.com/9seconds/guidedog/internal/options"
)

const envJSON = `{
	"hello": "world",
	"int_key": 1,
	"float_key": 1.1
}`

func createTempJSON(data string) string {
	fd, _ := ioutil.TempFile("", "")
	defer fd.Close()

	fd.WriteString(data)

	return fd.Name()
}

func createOptions() *opts.Options {
	options, _ := opts.NewOptions("term", // signal
		[]string{}, // envs
		0,          // gracefulTimeout
		"json",     // configFormat
		"",         // configPath
		[]string{}, // pathsToTracks
		"",         // lockFile
		false,      // pty
		false,      // supervise
		false)      // restartOnConfigChanges

	return options
}

func TestStringer(t *testing.T) {
	options := createOptions()
	env, err := NewEnvironment(options)

	assert.Nil(t, err)
	assert.True(t, env.String() != "")
}

func TestParseNoConfig(t *testing.T) {
	options := createOptions()
	env, _ := NewEnvironment(options)

	values, err := env.Parse()
	assert.Nil(t, err)
	assert.Equal(t, len(values), 0)
}

func TestParseWithConfig(t *testing.T) {
	options := createOptions()
	env, _ := NewEnvironment(options)

	configName := createTempJSON(envJSON)
	defer os.Remove(configName)

	env.Options.ConfigPath = configName

	values, err := env.Parse()

	assert.Nil(t, err)
	assert.Equal(t, len(values), 3)
	assert.Equal(t, values["hello"], "world")
	assert.Equal(t, values["int_key"], "1")
	assert.Equal(t, values["float_key"], "1.1")
}

func TestParseWithIncorrectConfig(t *testing.T) {
	options := createOptions()
	env, _ := NewEnvironment(options)

	configName := createTempJSON("")
	defer os.Remove(configName)

	env.Options.ConfigPath = configName

	_, err := env.Parse()

	assert.NotNil(t, err)
}

func TestUpdateWithConfigFormatNone(t *testing.T) {
	options := createOptions()
	env, _ := NewEnvironment(options)

	env.Options.ConfigFormat = opts.ConfigFormatNone

	err := env.Update()

	assert.Nil(t, err)
}

func TestUpdateWithIncorrectConfig(t *testing.T) {
	options := createOptions()
	env, _ := NewEnvironment(options)

	configName := createTempJSON("")
	defer os.Remove(configName)

	env.Options.ConfigPath = configName

	err := env.Update()

	assert.NotNil(t, err)
}

func TestUpdateWithProperConfig(t *testing.T) {
	options := createOptions()
	env, _ := NewEnvironment(options)

	configName := createTempJSON(envJSON)
	defer os.Remove(configName)

	env.Options.ConfigPath = configName

	err := env.Update()

	assert.Nil(t, err)
	assert.Equal(t, os.Getenv("hello"), "world")
	assert.Equal(t, os.Getenv("int_key"), "1")
	assert.Equal(t, os.Getenv("float_key"), "1.1")
}

func TestUpdateWithChangedConfig(t *testing.T) {
	options := createOptions()
	env, _ := NewEnvironment(options)

	configName := createTempJSON(envJSON)
	defer os.Remove(configName)

	env.Options.ConfigPath = configName

	env.Update()

	const changedJSON = "{\"hello\": \"bye\"}"
	ioutil.WriteFile(configName, []byte(changedJSON), os.FileMode(0666))

	err := env.Update()

	assert.Nil(t, err)
	assert.Equal(t, os.Getenv("hello"), "bye")
	assert.Equal(t, os.Getenv("int_key"), "")
	assert.Equal(t, os.Getenv("float_key"), "")
}

func TestUpdateWithPredefinedValues(t *testing.T) {
	options := createOptions()
	env, _ := NewEnvironment(options)

	configName := createTempJSON(envJSON)
	defer os.Remove(configName)

	env.Options.ConfigPath = configName
	env.Options.Envs = map[string]string{
		"int_key": "2",
	}

	env.Update()
	assert.Equal(t, os.Getenv("hello"), "world")
	assert.Equal(t, os.Getenv("int_key"), "2")
	assert.Equal(t, os.Getenv("float_key"), "1.1")

	const changedJSON = "{\"hello\": \"bye\"}"
	ioutil.WriteFile(configName, []byte(changedJSON), os.FileMode(0666))

	env.Update()
	assert.Equal(t, os.Getenv("hello"), "bye")
	assert.Equal(t, os.Getenv("int_key"), "2")
	assert.Equal(t, os.Getenv("float_key"), "")
}
