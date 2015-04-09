// Package execution contains all logic for execution of external commands
// based on Environment struct.
//
// This file contains supervising routines.
package execution

import (
	"fmt"
	"os"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

// supervisor defines structure which has all required data for supervising
// of running process.
type supervisor struct {
	allowedExitCodes  map[int]bool
	cmd               *command
	command           []string
	exitCodeChannel   chan int
	gracefulSignal    os.Signal
	gracefulTimeout   time.Duration
	hasTTY            bool
	keepAlivers       *sync.WaitGroup
	restartOnFailures bool
	supervisorChannel chan supervisorAction
}

func (s *supervisor) String() string {
	return fmt.Sprintf("%+v", *s)
}

// Just starts execution of the command and therefore its supervising.
func (s *supervisor) Start() {
	s.stop()

	if cmd, err := newCommand(s.command, s.hasTTY); err != nil {
		log.WithField("error", err).Panicf("Cannot start command!")
	} else {
		s.cmd = cmd
	}

	log.WithField("cmd", s.cmd).Info("Start process.")

	if s.restartOnFailures {
		go s.keepAlive()
	} else {
		go s.getExitCode()
	}
}

// Signal defines a callback for the incoming supervisorAction signal and
// reacts in expected way in a sync fashion.
func (s *supervisor) Signal(event supervisorAction) {
	switch event {
	case supervisorRestart:
		log.WithField("event", event).Info("Incoming restart event.")
		s.stop()
		s.Start()
	case supervisorStop:
		log.WithField("event", event).Info("Incoming stop event.")
		s.stop()
		s.exitCodeChannel <- s.cmd.ExitCode()
	}
}

// stopped just a thin wrapper which tells if command is stopped or not.
func (s *supervisor) stopped() bool {
	if s.cmd == nil {
		return true
	}
	return s.cmd.Stopped()
}

// stop just do what it names.
func (s *supervisor) stop() {
	log.Info("Stop external process.")

	log.Debug("Disable keepalivers.")
	s.keepAlivers.Wait()
	log.Debug("Keepalivers disabled.")

	if !s.stopped() {
		log.Debug("Start stopping process.")
		s.cmd.Stop(s.gracefulSignal, s.gracefulTimeout)
	} else {
		log.Debug("Process already stopped.")
	}
}

// keepAlive is just a function to be executed in goroutine. It tracks
// command execution and restarts if necessary.
func (s *supervisor) keepAlive() {
	s.keepAlivers.Add(1)
	defer s.keepAlivers.Done()

	log.Debug("Start keepaliver.")
	for {
		if s.stopped() {
			exitCode := s.cmd.ExitCode()
			if _, ok := s.allowedExitCodes[exitCode]; ok {
				log.WithFields(log.Fields{
					"exitCode":     exitCode,
					"allowedCodes": s.allowedExitCodes,
				}).Debug("Exit code means we have to stop the execution.")
				s.supervisorChannel <- supervisorStop
			} else {
				log.Debug("Process is stopped, restarting.")
				s.supervisorChannel <- supervisorRestart
			}

			log.Debug("Stop keepaliver.")
			return
		}
		time.Sleep(timeoutSupervising)
	}
}

// getExitCode has to be executed if no real supervising is performed. It just
// returns the command exit code.
func (s *supervisor) getExitCode() {
	for {
		if s.stopped() {
			s.exitCodeChannel <- s.cmd.ExitCode()
		}
		time.Sleep(timeoutSupervising)
	}
}

// newSupervisor returns new supervisor structure based on the given arguments.
// No command execution is performed at that moment.
func newSupervisor(command []string,
	exitCodeChannel chan int,
	gracefulSignal os.Signal,
	gracefulTimeout time.Duration,
	hasTTY bool,
	restartOnFailures bool,
	supervisorChannel chan supervisorAction,
	allowedExitCodes map[int]bool) *supervisor {
	return &supervisor{
		allowedExitCodes:  allowedExitCodes,
		command:           command,
		exitCodeChannel:   exitCodeChannel,
		gracefulSignal:    gracefulSignal,
		gracefulTimeout:   gracefulTimeout,
		hasTTY:            hasTTY,
		keepAlivers:       new(sync.WaitGroup),
		restartOnFailures: restartOnFailures,
		supervisorChannel: supervisorChannel,
	}
}
