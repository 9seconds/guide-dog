package execution

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

type SupervisorAction uint8

const (
	SUPERVISOR_START SupervisorAction = iota
	SUPERVISOR_STOP
	SUPERVISOR_RESTART
)

func NewSupervisorChannel() (channel chan SupervisorAction) {
	supervisorChannel := make(chan SupervisorAction, 1)

	go func() {
		signalChannel := makeSignalChannel()

		for signal := range signalChannel {
			log.WithField("signal", signal).Info("Signal received.")
			supervisorChannel <- SUPERVISOR_STOP
		}
	}()

	return supervisorChannel
}

func makeSignalChannel() (channel chan os.Signal) {
	channel = make(chan os.Signal, 1)

	signal.Notify(channel,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)

	return channel
}
