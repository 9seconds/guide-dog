// Package environment has a definition of Environment struct with parser.
// This file just defines Environment struct. For parsers please check
// parsers.go file.
package environment

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"

	opts "github.com/9seconds/guide-dog/options"
)

// Environment is just a thin container on opts.Options which can parse
// environment variables.
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

// Parse does parsing of the config file according to its ConfigFormat.
// Returns tuple of map with environment variables (key is the name, value
// is a, umm, value). Error defines the error.
func (env *Environment) Parse() (variables map[string]string, err error) {
	if env.Options.ConfigPath == "" {
		log.Info("Config path is not set, nothing to update.")
		return
	}

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

// Update does update of stored environment variables set with retrieved
// data from Parse output and maintains the set of environment variables
// of current process (which are derived by executed commands).
func (env *Environment) Update() (err error) {
	if env.Options.ConfigFormat == opts.ConfigFormatNone {
		return
	}

	variables, err := env.Parse()
	if err != nil {
		return
	}

	// Sets environment variables.
	for name, value := range variables {
		log.WithFields(log.Fields{
			"name":  name,
			"value": value,
		}).Debug("Set environment variable.")

		env.previousUpdates[name] = value
		os.Setenv(name, value)
	}

	// Maintains the list of previously set environment variables. Removes
	// obsoletes.
	for name := range env.previousUpdates {
		if _, ok := variables[name]; !ok {
			log.WithField("name", name).Debug("Delete environment variable.")
			delete(env.previousUpdates, name)
			os.Unsetenv(name)
		}
	}

	// Maintaines the list of explicitly preset environment variables.
	// Sets them forcefully, overrides previously set values.
	for name, value := range env.Options.Envs {
		log.WithFields(log.Fields{
			"name":  name,
			"value": value,
		}).Debug("Set predefined environment variable.")
		os.Setenv(name, value)
	}

	return
}

// NewEnvironment returns new Environment struct pointer and error if update
// failed.
func NewEnvironment(options *opts.Options) (env *Environment, err error) {
	env = &Environment{
		Options:         options,
		parser:          getParser(options),
		previousUpdates: make(map[string]string),
	}
	err = env.Update()

	return
}
