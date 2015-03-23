package execution

import (
	"time"
)

type SupervisorAction uint8

const (
	COMMAND_STILL_RUNNING       = -1
	COMMAND_INTERRUPT_EXIT_CODE = 130
	COMMAND_UNKNOWN_EXIT_CODE   = 70
)

const (
	GRACEFUL_SIGNAL_TIMEOUT = 2 * time.Millisecond
	SUPERVISOR_TIMEOUT      = 5 * time.Millisecond
	PTY_TIMEOUT             = 5 * time.Millisecond
	LOCK_FILE_TIMEOUT       = 5 * time.Millisecond
)

const (
	SUPERVISOR_STOP SupervisorAction = iota
	SUPERVISOR_RESTART
)
