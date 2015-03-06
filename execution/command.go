package execution

import (
	"os"
	"os/exec"
	"syscall"
	"time"
)

type Command struct {
	cmd             *exec.Cmd
	gracefulSignal  os.Signal
	gracefulTimeout time.Duration
}

func (c *Command) Stopped() bool {
	return c.cmd.ProcessState.Exited()
}

func (c *Command) ExitCode() int {
	waitStatus, ok := c.cmd.ProcessState.Sys().(syscall.WaitStatus)
	if !ok {
		return 1
	}

	return waitStatus.ExitStatus()
}

func (c *Command) Run() {
	c.cmd.Start()
}

func (c *Command) Stop() int {
	exitCodeChannel := make(chan int)
	gracefulTimerChannel := time.Tick(2 * time.Millisecond)
	killTimerChannel := time.After(c.gracefulTimeout)

	go func() {
		c.cmd.Process.Signal(c.gracefulSignal)
		defer c.cmd.Process.Release()

		for {
			select {
			case <-gracefulTimerChannel:
				if c.Stopped() {
					exitCodeChannel <- c.ExitCode()
					return
				}
			case <-killTimerChannel:
				if c.Stopped() {
					exitCodeChannel <- c.ExitCode()
				} else {
					c.cmd.Process.Kill()
				}
			}
		}
	}()

	return <-exitCodeChannel
}

func NewCommand(commandToExecute []string, gracefulSignal os.Signal,
	gracefulTimeout time.Duration, hasTTY bool) *Command {
	return &Command{
		cmd:             exec.Command(commandToExecute[0], commandToExecute[1:]...),
		gracefulSignal:  gracefulSignal,
		gracefulTimeout: gracefulTimeout,
	}
}
