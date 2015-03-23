package execution

import (
	"fmt"
	"os"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
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

func (sa SupervisorAction) String() string {
	switch sa {
	case SUPERVISOR_STOP:
		return "SUPERVISOR_STOP"
	case SUPERVISOR_RESTART:
		return "SUPERVISOR_RESTART"
	default:
		return "ERROR"
	}
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

	if cmd, err := NewCommand(s.command, s.hasTTY); err != nil {
		log.WithField("error", err).Panicf("Cannot start command!")
	} else {
		s.cmd = cmd
	}

	log.WithField("cmd", s.cmd).Info("Start process.")

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
		log.WithField("event", event).Info("Incoming restart event.")
		s.stop()
		s.Start()
	case SUPERVISOR_STOP:
		log.WithField("event", event).Info("Incoming stop event.")
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
	log.Info("Stop external process.")

	log.Debug("Disable keepalivers.")
	s.exit = true
	s.keepAlivers.Wait()
	log.Debug("Keepalivers disabled.")

	if !s.stopped() {
		log.Debug("Start stopping process.")
		s.cmd.Stop(s.gracefulSignal, s.gracefulTimeout)
	} else {
		log.Debug("Process already stopped.")
	}
}

func (s *Supervisor) keepAlive() {
	s.keepAlivers.Add(1)

	log.Debug("Start keepaliver.")
	for !s.exit {
		if s.stopped() {
			log.Debug("Process is stopped, restarting.")
			s.supervisorChannel <- SUPERVISOR_RESTART
		}
		time.Sleep(SUPERVISOR_TIMEOUT)
	}
	log.Debug("Stop keepaliver.")

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
