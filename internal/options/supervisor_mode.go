// Package options defines common options set for the guide-dog app.
package options

import "strings"

// SupervisorMode defines the mode of supervisor has to operate.
// Please check SupervisorMode* constants family for the possible
// values.
type SupervisorMode uint8

// SupervisorMode* consts family defines possible work modes of the supervisor,
// supported by the guide-dog.
const (
	SupervisorModeNone   SupervisorMode = 0
	SupervisorModeSimple SupervisorMode = 1 << iota
	SupervisorModeRestarting
)

func (sm SupervisorMode) String() string {
	mode := make([]string, 0, 2)

	if sm == SupervisorModeNone {
		mode = append(mode, "none")
	} else {
		if sm&SupervisorModeSimple > 0 {
			mode = append(mode, "simple")
		}
		if sm&SupervisorModeRestarting > 0 {
			mode = append(mode, "restarting")
		}
	}

	return strings.Join(mode, " / ")
}
