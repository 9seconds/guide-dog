package environment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func ConfigFormatJSONParser(filename string) (envs map[string]string, err error) {
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
