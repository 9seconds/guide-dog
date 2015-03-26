// Package options defines common options set for the guide-dog app.
package options

import (
	"fmt"
	"strings"
)

// ConfigFormat defines the type of the config on the given
// path. Please check ConfigFormat* constants family for the possible
// values.
type ConfigFormat uint8

// ConfigFormat* consts family defines possible config options, supported
// by the guide-dog.
const (
	ConfigFormatNone ConfigFormat = iota
	ConfigFormatJSON
	ConfigFormatYAML
	ConfigFormatINI
	ConfigFormatEnvDir
)

func (cf ConfigFormat) String() string {
	switch cf {
	case ConfigFormatNone:
		return "none"
	case ConfigFormatJSON:
		return "json"
	case ConfigFormatYAML:
		return "yaml"
	case ConfigFormatINI:
		return "ini"
	case ConfigFormatEnvDir:
		return "envdir"
	default:
		return "ERROR"
	}
}

func parseConfigFormat(name string) (format ConfigFormat, err error) {
	switch strings.ToLower(name) {
	case "":
		fallthrough
	case "none":
		format = ConfigFormatNone
	case "json":
		format = ConfigFormatJSON
	case "yaml":
		format = ConfigFormatYAML
	case "ini":
		format = ConfigFormatINI
	case "envdir":
		format = ConfigFormatEnvDir
	default:
		err = fmt.Errorf("Unknown config format %s", name)
	}

	return
}
