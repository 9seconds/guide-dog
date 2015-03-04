package main

import (
	"fmt"
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v1"

	environment "github.com/9seconds/guide-dog/environment"
	options "github.com/9seconds/guide-dog/options"
)

var (
	cmdLine = kingpin.New("guide-dog", "Small supervisor with envdir possibilities.")

	debug = cmdLine.
		Flag("debug", "Enable debug mode.").
		Short('d').
		Bool()
	signal = cmdLine.
		Flag("signal", "Signal to graceful shutting down of the given process.").
		Short('s').
		Default("SIGTERM").
		String()
	gracefulTimeout = cmdLine.
			Flag("graceful-tmo", "How long to wait for the process to be gracefully restarted. Before it got SIGKILLed.").
			Short('t').
			Default("5s").
			Duration()
	configFormat = cmdLine.
			Flag("config-format", "Format of configs.").
			Short('f').
			String()
	configPath = cmdLine.
			Flag("config-path", "Config path.").
			Short('p').
			String()
	lockFile = cmdLine.
			Flag("lock-file", "Lockfile on the local machine to acquire.").
			Short('l').
			String()
	supervise = cmdLine.
			Flag("supervise", "Set if it is required to supervise command. By default no supervising is performed.").
			Bool()
	superviseRestartOnConfigPathChanges = cmdLine.
						Flag("restart-on-config-changes", "Do the restart of the process if config is changed. Works only if 'supervise' option is enabled.").
						Bool()

	commandToExecute = cmdLine.
				Arg("command", "Command which has to be executed.").
				Strings()
)

func main() {
	kingpin.MustParse(cmdLine.Parse(os.Args[1:]))

	parsedOptions, err := options.NewOptions(*debug, *signal, *gracefulTimeout, *configFormat, *configPath, *lockFile, *supervise, *superviseRestartOnConfigPathChanges)
	if err != nil {
		panic(err)
	}

	env, err := environment.NewEnvironment(parsedOptions)

	fmt.Println(env)
}
