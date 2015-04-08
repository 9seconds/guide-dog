package main

import (
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	profile "github.com/davecheney/profile"
	kingpin "gopkg.in/alecthomas/kingpin.v1"

	environment "github.com/9seconds/guidedog/internal/environment"
	execution "github.com/9seconds/guidedog/internal/execution"
	options "github.com/9seconds/guidedog/internal/options"
)

const (
	profileEnvVariable = "GUIDEDOG_PROFILE"

	version = "0.1"

	envDirExitCode = 111
)

var (
	cmdLine = kingpin.New("guidedog", "Small supervisor with envdir possibilities.")

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
			Short('g').
			Default("SIGTERM").
			String()
	gracefulTimeout = cmdLine.
			Flag("graceful-tmo", "How long to wait for the process to be gracefully restarted. Before it got SIGKILLed.").
			Short('t').
			Default("5s").
			Duration()
	configFormat = cmdLine.
			Flag("config-format", "Format of configs.").
			Short('c').
			Enum("", "none", "json", "yaml", "ini", "envdir")
	configPath = cmdLine.
			Flag("config-path", "Config path.").
			Short('f').
			String()
	pathsToTrack = cmdLine.
			Flag("path-to-track", "Paths to track.").
			Short('p').
			Strings()
	lockFile = cmdLine.
			Flag("lock-file", "Lockfile on the local machine to acquire.").
			Short('l').
			String()
	pty = cmdLine.
		Flag("pty", "Allocate pseudo-terminal.").
		Short('y').
		Bool()
	runInShell = cmdLine.
			Flag("run-in-shell", "Run command in shell.").
			Short('x').
			Bool()
	supervise = cmdLine.
			Flag("supervise", "Set if it is required to supervise command. By default no supervising is performed.").
			Short('s').
			Bool()
	superviseRestartOnConfigPathChanges = cmdLine.
						Flag("restart-on-config-changes", "Do the restart of the process if config is changed. Works only if 'supervise' option is enabled.").
						Short('r').
						Bool()
	exitOnCodes = cmdLine.
			Flag("exit-on-code", "Exit if executed command finished with given code. You may define this option several times if you want to have several codes.").
			Short('o').
			Strings()
	commandToExecute = cmdLine.
				Arg("command", "Command which has to be executed.").
				Required().
				Strings()
)

// main is a classic entry point of any program.
func main() {
	exitCode := 0

	func() {
		defer func() {
			if exc := recover(); exc != nil {
				log.WithField("err", exc).Fatal("Fatal error happened.")
				exitCode = envDirExitCode
			}
		}()

		exitCode = mainWithExitCode()
	}()

	os.Exit(exitCode)
}

func init() {
	cmdLine.Version(version)
}

// mainWithExitCode returns the code to exit. This function is required
// because I want all deferred functions to be executed, os.Exit exits
// immediately. This is not cool.
func mainWithExitCode() int {
	kingpin.MustParse(cmdLine.Parse(os.Args[1:]))

	if os.Getenv(profileEnvVariable) != "" {
		defer profile.Start(profile.CPUProfile).Stop()
	}

	configureLogging(*debug)

	parsedOptions, err := options.NewOptions(
		*signalName,
		*envs,
		*gracefulTimeout,
		*configFormat,
		*configPath,
		*pathsToTrack,
		*lockFile,
		*pty,
		*supervise,
		*superviseRestartOnConfigPathChanges,
		*exitOnCodes)
	if err != nil {
		panic(err)
	}

	env, err := environment.NewEnvironment(parsedOptions)
	if err != nil {
		panic(err)
	}
	log.WithField("environment", env).Info("Environment.")

	if *runInShell {
		shell := os.Getenv("SHELL")
		*commandToExecute = []string{shell, "-i", "-c", strings.Join(*commandToExecute, " ")}
	}

	exitCode := execution.Execute(*commandToExecute, env)

	log.WithField("exitCode", exitCode).Info("Program exit")

	return exitCode
}

// configureLogging sets logging settings according to the debug option.
func configureLogging(debug bool) {
	log.SetOutput(os.Stderr)
	log.SetFormatter(&log.TextFormatter{})

	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.FatalLevel)
	}
}
