package environment

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	opts "github.com/9seconds/guide-dog/internal/options"

	assert "github.com/stretchr/testify/assert"
)

func getFileNameWithContent(data string) string {
	file, _ := ioutil.TempFile("", "")
	defer file.Close()

	file.WriteString(data)

	return file.Name()
}

func createEnvDirVariable(dir, name, content string) {
	path := filepath.Join(dir, name)
	ioutil.WriteFile(path, []byte(content), os.FileMode(0666))
}

func assertFuncEquals(t *testing.T, func1, func2 interface{}) {
	type1 := reflect.ValueOf(func1).Pointer()
	type2 := reflect.ValueOf(func2).Pointer()

	assert.Equal(t, type1, type2)
}

func TestGetParser(t *testing.T) {
	assertFuncEquals(t, getParser(opts.ConfigFormatNone), configFormatNoneParser)
	assertFuncEquals(t, getParser(opts.ConfigFormatJSON), configFormatJSONParser)
	assertFuncEquals(t, getParser(opts.ConfigFormatYAML), configFormatYAMLParser)
	assertFuncEquals(t, getParser(opts.ConfigFormatINI), configFormatINIParser)
	assertFuncEquals(t, getParser(opts.ConfigFormatEnvDir), configFormatEnvDirParser)
	assertFuncEquals(t, getParser(0xFF), configFormatNoneParser)
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

func TestConfigFormatINIOk(t *testing.T) {
	const data = `[somesection]
hello=world
foo = bar

[section2]
bar =   baaz`

	fileName := getFileNameWithContent(data)
	result, err := configFormatINIParser(fileName)

	assert.Nil(t, err)
	assert.Equal(t, len(result), 3)
	assert.Equal(t, result["hello"], "world")
	assert.Equal(t, result["foo"], "bar")
	assert.Equal(t, result["bar"], "baaz")
}

func TestConfigFormatINIFail(t *testing.T) {
	const data = `[somesec`

	fileName := getFileNameWithContent(data)
	_, err := configFormatINIParser(fileName)

	assert.NotNil(t, err)
}

func TestConfigFormatNoSections(t *testing.T) {
	const data = "t = 2"

	fileName := getFileNameWithContent(data)
	result, err := configFormatINIParser(fileName)

	assert.Nil(t, err)
	assert.Equal(t, len(result), 1)
	assert.Equal(t, result["t"], "2")
}

func TestConfigFormatEnvDirParserNoDir(t *testing.T) {
	_, err := configFormatEnvDirParser("WTF")

	assert.NotNil(t, err)
}

func TestConfigFormatEnvDirParserCorrectFiles(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tempDir)

	createEnvDirVariable(tempDir, "hello", "world")
	createEnvDirVariable(tempDir, "foo", "bar")

	result, err := configFormatEnvDirParser(tempDir)

	assert.Nil(t, err)
	assert.Equal(t, len(result), 2)
	assert.Equal(t, result["hello"], "world")
	assert.Equal(t, result["foo"], "bar")
}

func TestConfigFormatEnvDirParserEmptyFile(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tempDir)

	createEnvDirVariable(tempDir, "hello", "")

	result, err := configFormatEnvDirParser(tempDir)

	assert.Nil(t, err)
	assert.Equal(t, len(result), 1)
	assert.Equal(t, result["hello"], "")
}

func TestConfigFormatEnvDirParserSkipDirs(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tempDir)

	path := filepath.Join(tempDir, "WTF")
	os.Mkdir(path, os.FileMode(0666))

	result, err := configFormatEnvDirParser(tempDir)

	assert.Nil(t, err)
	assert.Equal(t, len(result), 0)
}
