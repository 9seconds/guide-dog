package environment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	ini "github.com/vaughan0/go-ini"

	yaml "gopkg.in/yaml.v2"
)

func configFormatNoneParser(path string) (envs map[string]string, err error) {
	return make(map[string]string), nil
}

func configFormatJSONParser(filename string) (envs map[string]string, err error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	var unmarshalled interface{}
	err = json.Unmarshal(content, &unmarshalled)
	if err != nil {
		return
	}

	convertedMap, converted := unmarshalled.(map[string]interface{})
	if !converted {
		return nil, fmt.Errorf("Incorrect JSON content in file %s", filename)
	}

	envs = make(map[string]string)
	for key, value := range convertedMap {
		strValue, converted := value.(string)
		if !converted {
			return nil, fmt.Errorf("Cannot convert %v to string", value)
		}
		envs[key] = strValue
	}

	return
}

func configFormatYAMLParser(filename string) (envs map[string]string, err error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	var unmarshalled interface{}
	err = yaml.Unmarshal(content, &unmarshalled)
	if err != nil {
		return
	}

	convertedMap, converted := unmarshalled.(map[string]interface{})
	if !converted {
		return nil, fmt.Errorf("Incorrect YAML content in file %s", filename)
	}

	envs = make(map[string]string)
	for key, value := range convertedMap {
		strValue, converted := value.(string)
		if !converted {
			return nil, fmt.Errorf("Cannot convert %v to string", value)
		}
		envs[key] = strValue
	}

	return
}

func configFormatINIParser(filename string) (envs map[string]string, err error) {
	file, err := ini.LoadFile(filename)
	if err != nil {
		return
	}

	envs = make(map[string]string)
	for _, data := range file {
		for key, value := range data {
			envs[key] = value
		}
	}

	return
}

func configFormatEnvDirParser(dirname string) (envs map[string]string, err error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return
	}

	envs = make(map[string]string)
	for _, item := range files {
		if item.IsDir() || item.Size() == 0 {
			continue
		}

		path := filepath.Join(dirname, item.Name())
		content, err := ioutil.ReadFile(path)
		if err != nil {
			continue
		}

		envs[item.Name()] = strings.TrimSpace(string(content))
	}

	return
}
