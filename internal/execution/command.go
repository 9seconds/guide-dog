// Package execution contains all logic for execution of external commands
// based on Environment struct.
//
// This file contains definition of command structure.
package execution

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	term "github.com/docker/docker/pkg/term"
	pty "github.com/kr/pty"
)

// command just a thin wrapper for the exec.Cmd which can restart
// and do some addition niceties.
type command struct {
	cmd *exec.Cmd
}

func (c *command) String() string {
	return fmt.Sprintf("%+v", c.cmd)
}

// Stopped checks if command stopped or not.
func (c *command) Stopped() bool {
	return c.cmd.ProcessState != nil
}

// ExitCode returns exit code of the command if it is stopped.
func (c *command) ExitCode() int {
	if !c.Stopped() {
		log.WithField("command", c).Warn("Command is still running!")
		return exitCodeStillRunning
	}

	waitStatus, ok := c.cmd.ProcessState.Sys().(syscall.WaitStatus)
	if !ok {
		log.Fatal("Cannot convert ProcessState to WaitStatus!")
		return exitCodeInternalError
	}

	exitCode := waitStatus.ExitStatus()
	if exitCode < 0 {
		exitCode = exitCodeInterrupt
	}

	return exitCode
}

// Stop do what the name defines.
func (c *command) Stop(signal os.Signal, timeout time.Duration) {
	if c.Stopped() {
		return
	}

	gracefulTimerChannel := time.Tick(timeoutGracefulSignal)
	killTimerChannel := time.After(timeout)

	log.WithField("cmd", c.cmd).Info("Start stopping process.")
	c.cmd.Process.Signal(signal)

	for !c.Stopped() {
		select {
		case <-gracefulTimerChannel:
			continue
		case <-killTimerChannel:
			log.Info("Graceful timeout expired, send kill signal")
			c.cmd.Process.Signal(syscall.SIGKILL)
		}
	}
}

// newCommand returns new running command instance.
func newCommand(commandToExecute []string, hasTTY bool) (commandToRun *command, err error) {
	cmd := exec.Command(commandToExecute[0], commandToExecute[1:]...)

	if hasTTY {
		cmd, err = makePTYCommand(cmd)
	} else {
		cmd, err = makeStandardCommand(cmd)
	}

	if err != nil {
		return
	}

	commandToRun = &command{cmd: cmd}

	go cmd.Wait()

	return
}

// makeStandardCommand just attach streams to the command and runs it.
func makeStandardCommand(cmd *exec.Cmd) (*exec.Cmd, error) {
	log.WithField("cmd", cmd).Info("Run command in standard mode.")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd, cmd.Start()
}

// makePTY command attaches streams to the command and run it with a
// preconfigured pseudo TTY.
func makePTYCommand(cmd *exec.Cmd) (*exec.Cmd, error) {
	log.WithField("cmd", cmd).Info("Run command with PTY.")

	pty, err := pty.Start(cmd)
	if err != nil {
		return cmd, err
	}

	hostFd := os.Stdin.Fd()
	oldTerminalState, err := term.SetRawTerminal(hostFd)
	if err != nil {
		return cmd, err
	}

	go func() {
		defer cleanUpPTY(cmd, pty, hostFd, oldTerminalState)

		for {
			if cmd.ProcessState != nil {
				return
			}
			time.Sleep(timeoutPTY)
		}
	}()

	monitorTTYResize(hostFd, pty.Fd())

	go io.Copy(pty, os.Stdin)
	go io.Copy(os.Stdout, pty)

	return cmd, nil
}

// cleanUpPTY closes configured PTY.
func cleanUpPTY(cmd *exec.Cmd, pty *os.File, hostFd uintptr, state *term.State) {
	log.WithField("cmd", cmd).Info("Cleanup PTY.")

	term.RestoreTerminal(hostFd, state)
	pty.Close()
}

// monitorTTYResize monitors if PTY winSize was changed and changes it
// in appropriate way.
func monitorTTYResize(hostFd uintptr, guestFd uintptr) {
	resizeTTY(hostFd, guestFd)

	winchChan := make(chan os.Signal, 1)
	signal.Notify(winchChan, syscall.SIGWINCH)

	go func() {
		for _ = range winchChan {
			resizeTTY(hostFd, guestFd)
		}
	}()
}

// resizeTTY just does what it names.
func resizeTTY(hostFd uintptr, guestFd uintptr) {
	winsize, err := term.GetWinsize(hostFd)

	if err != nil {
		return
	}

	term.SetWinsize(guestFd, winsize)
}
