package environment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	ini "github.com/vaughan0/go-ini"
	yaml "gopkg.in/yaml.v2"
)

type unmarshal func([]byte, interface{}) error

func configFormatNoneParser(path string) (envs map[string]string, err error) {
	return make(map[string]string), nil
}

func configUnmarshall(convertFromFloat bool, unpack unmarshal,
	filename string) (envs map[string]string, err error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Errorf("Cannot read from config file %s: %v", filename, err)
		return
	}

	var unmarshalled map[string]interface{}
	err = unpack(content, &unmarshalled)
	if err != nil {
		log.Errorf("Cannot unmarshal config file %s: %v", filename, err)
		return
	}
	log.Debugf("Unmarshalled structure %v", unmarshalled)

	envs = make(map[string]string)
	for key, value := range unmarshalled {
		if strValue, ok := value.(string); ok {
			envs[key] = strValue
			continue
		}

		if convertFromFloat {
			if floatValue, ok := value.(float64); ok {
				envs[key] = strconv.Itoa(int(floatValue))
				continue
			}
		}

		if intValue, ok := value.(int); ok {
			envs[key] = strconv.Itoa(intValue)
			continue
		}

		log.Errorf("Cannot convert %v to string", value)
		return nil, fmt.Errorf("Cannot convert %v to string", value)
	}

	return
}

func configFormatJSONParser(filename string) (map[string]string, error) {
	return configUnmarshall(true, json.Unmarshal, filename)
}

func configFormatYAMLParser(filename string) (map[string]string, error) {
	return configUnmarshall(false, yaml.Unmarshal, filename)
}

func configFormatINIParser(filename string) (envs map[string]string, err error) {
	file, err := ini.LoadFile(filename)
	if err != nil {
		log.Errorf("Cannot read from config file %s: %v", filename, err)
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
		log.Errorf("Cannot read directory %s: %v", dirname, err)
		return
	}

	envs = make(map[string]string)
	for _, item := range files {
		if item.IsDir() {
			log.Debugf("%s is a directory, skip", item.Name())
			continue
		}

		if item.Size() == 0 {
			log.Debugf("%s has 0 size, set %s to empty string", item.Name(), item.Name())
			envs[item.Name()] = ""
			continue
		}

		path := filepath.Join(dirname, item.Name())
		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Warnf("Cannot read file %s, skip: %v", path, err)
			continue
		}

		envs[item.Name()] = strings.TrimSpace(string(content))
	}

	return
}
