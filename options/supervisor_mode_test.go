package options

import (
	"strings"
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestSupervisorModeName(t *testing.T) {
	assert.Equal(t, "none", SupervisorModeNone.String())
	assert.Equal(t, "simple", SupervisorModeSimple.String())
	assert.Equal(t, "restarting", SupervisorModeRestarting.String())

	mode := SupervisorModeSimple | SupervisorModeRestarting
	assert.True(t, strings.Contains(mode.String(), "simple"))
	assert.True(t, strings.Contains(mode.String(), "restarting"))
	assert.True(t, !strings.Contains(mode.String(), "none"))
}
