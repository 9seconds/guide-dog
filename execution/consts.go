// Package execution contains all logic for execution of external commands
// based on Environment struct.
//
// This file contains common constants.
package execution

import (
	"time"
)

// supervisorAction defines the action which is required to be performed
// with executing command (restart or stop).
type supervisorAction uint8

// exitCode* constants family defines exit codes for managed situations.
const (
	exitCodeStillRunning  = -1
	exitCodeInterrupt     = 130
	exitCodeInternalError = 70
)

// timeout* constants family defines time.Durations for different
// internal purposes.
const (
	timeoutGracefulSignal = 2 * time.Millisecond
	timeoutSupervising    = 5 * time.Millisecond
	timeoutPTY            = 5 * time.Millisecond
	timeoutLockFile       = 5 * time.Millisecond
)

// supervisor* constants family defines the set of actions that could be
// performed during process supervising.
const (
	supervisorStop supervisorAction = iota
	supervisorRestart
)

func (sa supervisorAction) String() string {
	switch sa {
	case supervisorStop:
		return "SupervisorStop"
	case supervisorRestart:
		return "SupervisorRestart"
	default:
		return "ERROR"
	}
}
