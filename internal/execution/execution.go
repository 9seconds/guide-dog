// Package execution contains all logic for execution of external commands
// based on Environment struct.
//
// This file contains Execute function.
package execution

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"

	environment "github.com/9seconds/guidedog/internal/environment"
	options "github.com/9seconds/guidedog/internal/options"
)

// Execute just executes given command in with given Environment.
// It configures supervising if necessary, filesystem notifications etc.
// It does work.
func Execute(command []string, env *environment.Environment) int {
	if env.Options.LockFile != nil {
		for {
			if err := env.Options.LockFile.Acquire(); err == nil {
				defer env.Options.LockFile.Release()
				break
			}
			time.Sleep(timeoutLockFile)
		}
	}

	pathsToWatch := []string{env.Options.ConfigPath}
	pathsToWatch = append(pathsToWatch, env.Options.PathsToTrack...)

	watcherChannel := makeWatcher(pathsToWatch, env)
	defer close(watcherChannel)

	exitCodeChannel := make(chan int, 1)
	defer close(exitCodeChannel)

	supervisorChannel := make(chan supervisorAction, 1)
	defer close(supervisorChannel)

	signalChannel := makeSignalChannel()
	defer close(signalChannel)

	go attachSignalChannel(supervisorChannel, signalChannel)
	if env.Options.Supervisor&options.SupervisorModeRestarting > 0 {
		go attachSupervisorChannel(supervisorChannel, watcherChannel)
	}

	supervisor := newSupervisor(command,
		exitCodeChannel,
		env.Options.Signal,
		env.Options.GracefulTimeout,
		env.Options.PTY,
		env.Options.Supervisor&options.SupervisorModeSimple > 0,
		supervisorChannel,
		env.Options.ExitCodes)

	log.WithField("supervisor", supervisor).Info("Start supervisor.")

	supervisor.Start()
	go func() {
		for {
			event, ok := <-supervisorChannel
			if !ok {
				return
			}
			supervisor.Signal(event)
		}
	}()

	return <-exitCodeChannel
}

// attachSignalChannel attaches given signalChannel events and configures
// basic supervising actions. Basically it just sends stop signal to external
// command on interrupt.
func attachSignalChannel(channel chan supervisorAction, signalChannel chan os.Signal) {
	for {
		incomingSignal, ok := <-signalChannel
		if !ok {
			return
		}

		log.WithField("signal", incomingSignal).Debug("Signal from OS received.")

		channel <- supervisorStop
	}
}

// attachSupervisorChannel attaches some restart supervisor channel (filesystem
// notifications for example) to common supervisorAction channel.
func attachSupervisorChannel(channel chan supervisorAction, supervisorChannel chan bool) {
	for {
		event, ok := <-supervisorChannel
		if !ok {
			return
		}

		log.WithFields(log.Fields{
			"event":   event,
			"channel": supervisorChannel,
		}).Debug("Event from supervisor channel is captured.")

		channel <- supervisorRestart
	}
}

// makeSignalChannel is a generic routine which connects signal handler
// to the channel.
func makeSignalChannel() (channel chan os.Signal) {
	channel = make(chan os.Signal, 1)

	signal.Notify(channel,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)

	return channel
}
