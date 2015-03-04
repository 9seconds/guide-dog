package environment

import (
	opts "github.com/9seconds/guide-dog/options"
)

type environmentParser func(string) (map[string]string, error)

type Environment struct {
	Options *opts.Options

	parser environmentParser
}

func (env *Environment) Parse() (map[string]string, error) {
	return env.parser(env.Options.ConfigPath)
}

func NewEnvironment(options *opts.Options) (env *Environment, err error) {
	env = &Environment{
		Options: options,
		parser:  getParser(options),
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
