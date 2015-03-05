package environment

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"

	opts "github.com/9seconds/guide-dog/options"
)

type environmentParser func(string) (map[string]string, error)

type Environment struct {
	Options         *opts.Options
	parser          environmentParser
	previousUpdates map[string]string
}

func (env *Environment) String() string {
	return fmt.Sprintf("<Environment(Options='%v', parser='%v', previousUpdates='%v')>",
		env.Options,
		env.parser,
		env.previousUpdates,
	)
}

func (env *Environment) Parse() (variables map[string]string, err error) {
	variables, err = env.parser(env.Options.ConfigPath)
	if err != nil {
		log.WithFields(log.Fields{
			"configPath": env.Options.ConfigPath,
			"error":      err,
		}).Warn("Cannot parse")
	} else {
		log.WithField("variables", variables).Info("Parsed environment variables.")
	}

	return
}

func (env *Environment) Update() (err error) {
	variables, err := env.Parse()
	if err != nil {
		return
	}

	for name, value := range variables {
		log.WithFields(log.Fields{
			"name":  name,
			"value": value,
		}).Debug("Set environment variable.")

		env.previousUpdates[name] = value
		os.Setenv(name, value)
	}

	for name, _ := range env.previousUpdates {
		if _, ok := variables[name]; !ok {
			log.WithField("name", name).Debug("Delete environment variable.")
			delete(env.previousUpdates, name)
			os.Unsetenv(name)
		}
	}

	for name, value := range env.Options.Envs {
		log.WithFields(log.Fields{
			"name":  name,
			"value": value,
		}).Debug("Set predefined environment variable.")
		os.Setenv(name, value)
	}

	return
}

func NewEnvironment(options *opts.Options) (env *Environment, err error) {
	env = &Environment{
		Options:         options,
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
