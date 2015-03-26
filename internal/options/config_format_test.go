package options

import (
	"strings"
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestParseConfigFormat(t *testing.T) {
	validNames := []string{"", "none", "json", "yaml", "ini", "envdir"}
	formats := []ConfigFormat{ConfigFormatNone, ConfigFormatNone,
		ConfigFormatJSON, ConfigFormatYAML, ConfigFormatINI, ConfigFormatEnvDir}

	for idx, name := range validNames {
		for _, caseSensitiveName := range []string{name, strings.ToUpper(name)} {
			format, err := parseConfigFormat(caseSensitiveName)
			assert.Nil(t, err)
			assert.Equal(t, formats[idx], format)
		}
	}
}

func TestParseUnknownConfigFormat(t *testing.T) {
	_, err := parseConfigFormat("WTF")
	assert.NotNil(t, err)
}

func TestConfigFormatNames(t *testing.T) {
	assert.Equal(t, ConfigFormatNone.String(), "none")
	assert.Equal(t, ConfigFormatJSON.String(), "json")
	assert.Equal(t, ConfigFormatYAML.String(), "yaml")
	assert.Equal(t, ConfigFormatINI.String(), "ini")
	assert.Equal(t, ConfigFormatEnvDir.String(), "envdir")
}
