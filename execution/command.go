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

type Command struct {
	cmd *exec.Cmd
}

func (c *Command) String() string {
	return fmt.Sprintf("<Command(command='%v' (%v))>",
		c.cmd.Args,
		c.cmd)
}

func (c *Command) Stopped() bool {
	return c.cmd.ProcessState != nil
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

	exitCode := waitStatus.ExitStatus()
	if exitCode < 0 {
		exitCode = COMMAND_INTERRUPT_EXIT_CODE
	}

	return exitCode
}

func (c *Command) Stop(signal os.Signal, timeout time.Duration) {
	if c.Stopped() {
		return
	}

	gracefulTimerChannel := time.Tick(GRACEFUL_SIGNAL_TIMEOUT)
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

func NewCommand(commandToExecute []string, hasTTY bool) (command *Command, err error) {
	cmd := exec.Command(commandToExecute[0], commandToExecute[1:]...)

	if hasTTY {
		cmd, err = makePTYCommand(cmd)
	} else {
		cmd, err = makeStandardCommand(cmd)
	}

	if err != nil {
		return
	}

	command = &Command{cmd: cmd}

	go cmd.Wait()

	return
}

func makeStandardCommand(cmd *exec.Cmd) (*exec.Cmd, error) {
	log.WithField("cmd", cmd).Info("Run command in standard mode.")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd, cmd.Start()
}

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
		defer cleanupPTY(cmd, pty, hostFd, oldTerminalState)

		for {
			if cmd.ProcessState != nil {
				return
			}
			time.Sleep(PTY_TIMEOUT)
		}
	}()

	monitorTTYResize(hostFd, pty.Fd())

	go io.Copy(pty, os.Stdin)
	go io.Copy(os.Stdout, pty)

	return cmd, nil
}

func cleanupPTY(cmd *exec.Cmd, pty *os.File, hostFd uintptr, state *term.State) {
	log.WithField("cmd", cmd).Info("Cleanup PTY.")

	term.RestoreTerminal(hostFd, state)
	pty.Close()
}

func monitorTTYResize(hostFd uintptr, guestFd uintptr) {
	resizeTty(hostFd, guestFd)

	winchChan := make(chan os.Signal, 1)
	signal.Notify(winchChan, syscall.SIGWINCH)

	go func() {
		for _ = range winchChan {
			resizeTty(hostFd, guestFd)
		}
	}()
}

func resizeTty(hostFd uintptr, guestFd uintptr) {
	winsize, err := term.GetWinsize(hostFd)

	if err != nil {
		return
	}

	term.SetWinsize(guestFd, winsize)
}
