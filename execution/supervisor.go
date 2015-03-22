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
	}
}

func (s *Supervisor) keepAlive() {
	s.keepAlivers.Add(1)
	defer s.keepAlivers.Done()

	log.Debug("Start keepAlive.")
	for !s.exit {
		if s.stopped() {
			log.Debug("Process is stopped, restarting.")
			s.supervisorChannel <- SUPERVISOR_RESTART
		}
		time.Sleep(20 * time.Millisecond)
	}
	log.Debug("Finish keepAlive.")
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
	return s.cmd != nil && s.cmd.Stopped()
}

func (s *Supervisor) stop() {
	if !s.stopped() {
		s.exit = true
		s.keepAlivers.Wait()
		s.cmd.Stop(s.gracefulSignal, s.gracefulTimeout)
	}
}

func NewSupervisor(command []string, exitCodeChannel chan int, gracefulSignal os.Signal, gracefulTimeout time.Duration, hasTTY bool, restartOnFailures bool, supervisorChannel chan SupervisorAction) *Supervisor {
	return &Supervisor{
		command:           command,
		exit:              false,
		exitCodeChannel:   exitCodeChannel,
		gracefulSignal:    gracefulSignal,
		gracefulTimeout:   gracefulTimeout,
		hasTTY:            hasTTY,
		supervisorChannel: supervisorChannel,
		keepAlivers:       new(sync.WaitGroup)}
}
