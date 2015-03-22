package execution

import (
	"fmt"
	"os"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

type SupervisorAction uint8

const (
	SUPERVISOR_STOP SupervisorAction = iota
	SUPERVISOR_RESTART

	SUPERVISOR_TIMEOUT = 5 * time.Millisecond
)

type Supervisor struct {
	cmd               *Command
	command           []string
	exit              bool
	exitCodeChannel   chan int
	gracefulSignal    os.Signal
	gracefulTimeout   time.Duration
	hasTTY            bool
	restartOnFailures bool
	supervisorChannel chan SupervisorAction
	keepAlivers       *sync.WaitGroup
}

func (s *Supervisor) String() string {
	return fmt.Sprintf("<Supervisor(cmd='%v', command='%v', exit=%t, exitCodeChannel='%v', gracefulSignal=%v, gracefulTimeout=%v, hasTTY=%t, restartOnFailures=%t, supervisorChannel=%v, keepAlivers=%v)>",
		s.cmd,
		s.command,
		s.exit,
		s.exitCodeChannel,
		s.gracefulSignal,
		s.gracefulTimeout,
		s.hasTTY,
		s.restartOnFailures,
		s.supervisorChannel,
		s.keepAlivers)
}

func (s *Supervisor) Start() {
	s.stop()

	s.cmd = NewCommand(s.command, s.hasTTY)
	s.cmd.Start()

	s.exit = false
	if s.restartOnFailures {
		go s.keepAlive()
	} else {
		go s.getExitCode()
	}
}

func (s *Supervisor) Signal(event SupervisorAction) {
	switch event {
	case SUPERVISOR_RESTART:
		s.stop()
		s.Start()
	case SUPERVISOR_STOP:
		s.stop()
		s.exitCodeChannel <- s.cmd.ExitCode()
	}
}

func (s *Supervisor) stopped() bool {
	if s.cmd == nil {
		return true
	}
	return s.cmd.Stopped()
}

func (s *Supervisor) stop() {
	s.exit = true
	s.keepAlivers.Wait()

	if !s.stopped() {
		s.cmd.Stop(s.gracefulSignal, s.gracefulTimeout)
	}
}

func (s *Supervisor) keepAlive() {
	s.keepAlivers.Add(1)

	for !s.exit {
		if s.stopped() {
			log.Debug("Process is stopped, restarting.")
			s.supervisorChannel <- SUPERVISOR_RESTART
		}
		time.Sleep(SUPERVISOR_TIMEOUT)
	}

	s.keepAlivers.Done()
}

func (s *Supervisor) getExitCode() {
	for {
		if s.stopped() {
			s.exitCodeChannel <- s.cmd.ExitCode()
		}
		time.Sleep(SUPERVISOR_TIMEOUT)
	}
}

func NewSupervisor(command []string, exitCodeChannel chan int, gracefulSignal os.Signal, gracefulTimeout time.Duration, hasTTY bool, restartOnFailures bool, supervisorChannel chan SupervisorAction) *Supervisor {
	return &Supervisor{
		command:           command,
		exitCodeChannel:   exitCodeChannel,
		exit:              false,
		gracefulSignal:    gracefulSignal,
		gracefulTimeout:   gracefulTimeout,
		hasTTY:            hasTTY,
		keepAlivers:       new(sync.WaitGroup),
		restartOnFailures: restartOnFailures,
		supervisorChannel: supervisorChannel}
}
