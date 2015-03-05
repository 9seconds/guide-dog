package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v1"
	fsnotify "gopkg.in/fsnotify.v1"

	environment "github.com/9seconds/guide-dog/environment"
	options "github.com/9seconds/guide-dog/options"
)

const ENV_DIR_EXIT_CODE = 111

var (
	cmdLine = kingpin.New("guide-dog", "Small supervisor with envdir possibilities.")

	debug = cmdLine.
		Flag("debug", "Enable debug mode.").
		Short('d').
		Bool()
	envs = cmdLine.
		Flag("env", "Environment variable to set. There may be several options '-e OS=Linux -e H=1'.").
		Short('e').
		Strings()
	signalName = cmdLine.
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
	pty = cmdLine.
		Flag("pty", "Allocate pseudo-terminal.").
		Bool()
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
	defer func() {
		if exc := recover(); exc != nil {
			log.WithField("err", exc).Fatal("Fatal error happened.")
			os.Exit(ENV_DIR_EXIT_CODE)
		}
	}()

	kingpin.MustParse(cmdLine.Parse(os.Args[1:]))

	configureLogging(*debug)

	parsedOptions, err := options.NewOptions(
		*debug,
		*signalName,
		*envs,
		*gracefulTimeout,
		*configFormat,
		*configPath,
		*lockFile,
		*pty,
		*supervise,
		*superviseRestartOnConfigPathChanges,
	)
	if err != nil {
		panic(err)
	}

	env, err := environment.NewEnvironment(parsedOptions)
	if err != nil {
		panic(err)
	}
	log.WithField("environment", env).Info("Environment.")

	if len(*commandToExecute) > 0 {
		exitCode := execute(*commandToExecute, env)
		log.WithField("exitCode", exitCode).Info("Program exit")
		os.Exit(exitCode)
	}
}

func configureLogging(debug bool) {
	log.SetOutput(os.Stderr)
	log.SetFormatter(&log.TextFormatter{})

	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.FatalLevel)
	}
}

func execute(commandToExecute []string, env *environment.Environment) int {
	exitCodeChannel := make(chan int, 1)

	watcher, watcherChannel := makeWatcher(env.Options.ConfigPath)
	defer close(watcherChannel)
	defer watcher.Close()

	log.Info(len(watcherChannel))

	<-exitCodeChannel

	return 0
}

func makeWatcher(configPath string) (watcher *fsnotify.Watcher, channel chan bool) {
	channel = make(chan bool, 1)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	if configPath == "" {
		return
	}

	err = watcher.Add(configPath)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op == 0 {
					continue
				}

				log.WithFields(log.Fields{
					"event": event,
					"op":    event.Op,
				}).Info("Event from filesystem is coming")

				if len(channel) == 0 {
					channel <- true
				}
			case err := <-watcher.Errors:
				if err != nil {
					log.WithField("error", err).Error("Some problem with filesystem notifications")
				}
			}
		}
	}()

	return
}

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
