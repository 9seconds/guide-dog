package environment

import (
	"fmt"
	"os"

	opts "github.com/9seconds/guide-dog/options"
	log "github.com/Sirupsen/logrus"
)

type environmentParser func(string) (map[string]string, error)

type Environment struct {
	options         *opts.Options
	parser          environmentParser
	previousUpdates map[string]string
}

func (env *Environment) String() string {
	return fmt.Sprintf("<Environment(options='%v', parser='%v', previousUpdates='%v')>",
		env.options,
		env.parser,
		env.previousUpdates,
	)
}

func (env *Environment) Parse() (variables map[string]string, err error) {
	variables, err = env.parser(env.options.ConfigPath)
	if err != nil {
		log.Warnf("Cannot parse %s: %v", env.options.ConfigPath, err)
	} else {
		log.Debugf("Parsed environment variables are %v", variables)
	}

	return
}

func (env *Environment) Update() (err error) {
	variables, err := env.Parse()
	if err != nil {
		log.Warnf("Cannot parse environment variables, skip: %v", err)
		return
	}

	for name, value := range variables {
		log.Debugf("Set environment variable %s to %s", name, value)

		env.previousUpdates[name] = value
		os.Setenv(name, value)
	}

	for name, _ := range env.previousUpdates {
		if _, ok := variables[name]; !ok {
			log.Debugf("Delete environment variable %s", name)
			delete(env.previousUpdates, name)
			// TODO: go 1.4
			// os.Unsetenv(name)
		}
	}

	for name, value := range env.options.Envs {
		log.Debugf("Set predefined environment variable %s to %s", name, value)
		os.Setenv(name, value)
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
