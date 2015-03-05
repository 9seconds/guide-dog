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
		log.WithFields(log.Fields{
			"filename": filename,
			"error":    err,
		}).Error("Cannot read from config file.")
		return
	}

	var unmarshalled map[string]interface{}
	err = unpack(content, &unmarshalled)
	if err != nil {
		log.WithFields(log.Fields{
			"filename": filename,
			"error":    err,
		}).Error("Cannot parse config file")
		return
	}
	log.WithField("structure", unmarshalled).Debug("Unmarshalled structure.")

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

		log.WithField("value", value).Error("Cannot convert to string.")
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
		log.WithFields(log.Fields{
			"filename": filename,
			"error":    err,
		}).Error("Cannot read from config file.")
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
		log.WithFields(log.Fields{
			"path":  dirname,
			"error": err,
		}).Error("Cannot list directory.")
		return
	}

	envs = make(map[string]string)
	for _, item := range files {
		if item.IsDir() {
			log.WithFields(log.Fields{
				"dirname": dirname,
				"name":    item.Name(),
			}).Debug("Skip directory.")
			continue
		}

		if item.Size() == 0 {
			log.WithFields(log.Fields{
				"dirname": dirname,
				"name":    item.Name(),
			}).Debug("Set to empty string because filesize is 0.")
			envs[item.Name()] = ""
			continue
		}

		path := filepath.Join(dirname, item.Name())
		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.WithFields(log.Fields{
				"dirname": dirname,
				"name":    item.Name(),
			}).Warn("Cannot read file, skip.")
			continue
		}

		envs[item.Name()] = strings.TrimSpace(string(content))
	}

	return
}
