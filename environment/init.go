package environment

import (
	"os"

	opts "github.com/9seconds/guide-dog/options"
)

type environmentParser func(string) (map[string]string, error)

type Environment struct {
	options         *opts.Options
	parser          environmentParser
	previousUpdates map[string]string
}

func (env *Environment) Parse() (map[string]string, error) {
	return env.parser(env.options.ConfigPath)
}

func (env *Environment) Update() (err error) {
	variables, err := env.Parse()
	if err != nil {
		return
	}

	for name, value := range variables {
		env.previousUpdates[name] = value
		os.Setenv(name, value)
	}
	for name, _ := range env.previousUpdates {
		if _, ok := variables[name]; !ok {
			delete(env.previousUpdates, name)
			// TODO: go 1.4
			// os.Unsetenv(name)
		}
	}

	return
}

func NewEnvironment(options *opts.Options) (env *Environment, err error) {
	env = &Environment{
		options:         options,
		parser:          getParser(options),
		previousUpdates: make(map[string]string),
	}

	return
}

func getParser(options *opts.Options) environmentParser {
	switch options.ConfigFormat {
	case opts.CONFIG_FORMAT_NONE:
		return configFormatNoneParser
	case opts.CONFIG_FORMAT_JSON:
		return configFormatJSONParser
	case opts.CONFIG_FORMAT_YAML:
		return configFormatYAMLParser
	case opts.CONFIG_FORMAT_INI:
		return configFormatINIParser
	case opts.CONFIG_FORMAT_ENVDIR:
		return configFormatEnvDirParser
	default:
		return configFormatNoneParser
	}
}
