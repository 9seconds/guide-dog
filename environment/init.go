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
		return ConfigFormatNoneParser
	case opts.CONFIG_FORMAT_JSON:
		return ConfigFormatJSONParser
	// case opts.CONFIG_FORMAT_YAML:
	// 	return ConfigFormatYAMLParser
	// case opts.CONFIG_FORMAT_INI:
	// 	return ConfigFormatINIParser
	// case opts.CONFIG_FORMAT_ENVDIR:
	// 	return ConfigFormatEnvDirParser
	default:
		return ConfigFormatNoneParser
	}
}
