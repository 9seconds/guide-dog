package execution

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
)

type Command struct {
	cmd             *exec.Cmd
	gracefulSignal  os.Signal
	gracefulTimeout time.Duration
}

const (
	COMMAND_STILL_RUNNING     = -1
	COMMAND_UNKNOWN_EXIT_CODE = 70
)

func (c *Command) String() string {
	return fmt.Sprintf("<Command(command='%v' (%v), gracefulSignal='%v', gracefulTimeout='%v')>",
		c.cmd.Args,
		c.cmd,
		c.gracefulSignal,
		c.gracefulTimeout,
	)
}

func (c *Command) Stopped() bool {
	return c.cmd.ProcessState != nil && c.cmd.ProcessState.Exited()

}

func (c *Command) ExitCode() int {
	if !c.Stopped() {
		log.WithField("command", c).Warn("Command is still running!")
		return COMMAND_STILL_RUNNING
	}

	waitStatus, ok := c.cmd.ProcessState.Sys().(syscall.WaitStatus)
	if !ok {
		log.Fatal("Cannot convert ProcessState to WaitStatus!")
		return COMMAND_UNKNOWN_EXIT_CODE
	}

	return waitStatus.ExitStatus()
}

func (c *Command) Run() {
	c.cmd.Start()
}

func (c *Command) Stop() int {
	gracefulTimerChannel := time.Tick(2 * time.Millisecond)
	killTimerChannel := time.After(c.gracefulTimeout)

	c.cmd.Process.Signal(c.gracefulSignal)
	defer c.cmd.Process.Release()

	for {
		select {
		case <-gracefulTimerChannel:
			if c.Stopped() {
				break
			}
		case <-killTimerChannel:
			if c.Stopped() {
				break
			} else {
				log.Info("Graceful timeout expired, send kill signal")
				c.cmd.Process.Kill()
			}
		}
	}

	return c.ExitCode()
}

func NewCommand(commandToExecute []string, gracefulSignal os.Signal,
	gracefulTimeout time.Duration, hasTTY bool) *Command {
	return &Command{
		cmd:             exec.Command(commandToExecute[0], commandToExecute[1:]...),
		gracefulSignal:  gracefulSignal,
		gracefulTimeout: gracefulTimeout,
	}
}
