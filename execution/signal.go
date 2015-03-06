package execution

import (
	"os"
	"os/signal"
	"syscall"
)

func makeSignalChannel() (channel chan os.Signal) {
	channel = make(chan os.Signal, 1)

	signal.Notify(channel,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTSTP,
		syscall.SIGCONT,
		syscall.SIGTTIN,
		syscall.SIGTTOU,
		syscall.SIGBUS,
		syscall.SIGSEGV,
	)

	return channel
}
