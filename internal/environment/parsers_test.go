package environment

import (
	"io/ioutil"
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func getFileNameWithContent(data string) string {
	file, _ := ioutil.TempFile("", "")
	defer file.Close()

	file.WriteString(data)

	return file.Name()
}

func TestConfigFormatJSONOk(t *testing.T) {
	const data = `
	{
		"hello": "world",
		"key": "value",
		"int": 1,
		"float_ceiled": 1.0,
		"float_as_float": 1.0001
	}
	`

	fileName := getFileNameWithContent(data)
	result, err := configFormatJSONParser(fileName)

	assert.Nil(t, err)
	assert.Equal(t, len(result), 5)
	assert.Equal(t, result["hello"], "world")
	assert.Equal(t, result["key"], "value")
	assert.Equal(t, result["int"], "1")
	assert.Equal(t, result["float_ceiled"], "1")
	assert.Equal(t, result["float_as_float"], "1.0001")
}

func TestConfigFormatJSONEmpty(t *testing.T) {
	fileName := getFileNameWithContent("")
	_, err := configFormatJSONParser(fileName)

	assert.NotNil(t, err)
}

func TestConfigFormatJSONCorrupted(t *testing.T) {
	const data = `
	{
		"
	}
	`

	fileName := getFileNameWithContent(data)
	_, err := configFormatJSONParser(fileName)

	assert.NotNil(t, err)
}

func TestConfigFormatJSONIncorrectValue(t *testing.T) {
	const data = `
	{
		"hello": []
	}
	`

	fileName := getFileNameWithContent(data)
	_, err := configFormatJSONParser(fileName)

	assert.NotNil(t, err)
}

func TestConfigFormatYAMLOk(t *testing.T) {
	const data = `hello: world
key: value
int: 1
float_ceiled: 1.0
float_as_float: 1.0001
`

	fileName := getFileNameWithContent(data)
	result, err := configFormatYAMLParser(fileName)

	assert.Nil(t, err)
	assert.Equal(t, len(result), 5)
	assert.Equal(t, result["hello"], "world")
	assert.Equal(t, result["key"], "value")
	assert.Equal(t, result["int"], "1")
	assert.Equal(t, result["float_ceiled"], "1")
	assert.Equal(t, result["float_as_float"], "1.0001")
}

func TestConfigFormatYAMLEmpty(t *testing.T) {
	fileName := getFileNameWithContent("")
	_, err := configFormatYAMLParser(fileName)

	assert.Nil(t, err)
}

func TestConfigFormatYAMLCorrupted(t *testing.T) {
	const data = `
	{
		"
	}
	`

	fileName := getFileNameWithContent(data)
	_, err := configFormatYAMLParser(fileName)

	assert.NotNil(t, err)
}

func TestConfigFormatYAMLIncorrectValue(t *testing.T) {
	const data = `hello:
	- world
toot-toot:
	lala: 1
`

	fileName := getFileNameWithContent(data)
	_, err := configFormatJSONParser(fileName)

	assert.NotNil(t, err)
}
