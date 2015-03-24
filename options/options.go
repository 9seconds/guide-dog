// Package options defines common options set for the guide-dog app.
package options

import (
	"fmt"
	"strings"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"

	lockfile "github.com/9seconds/guide-dog/lockfile"
)

type (
	// ConfigFormat defines the type of the config on the given
	// path. Please check ConfigFormat* constants family for the possible
	// values.
	ConfigFormat uint8

	// SupervisorMode defines the mode of supervisor has to operate.
	// Please check SupervisorMode* constants family for the possible
	// values.
	SupervisorMode uint8
)

// ConfigFormat* consts family defines possible config options, supported
// by the guide-dog.
const (
	ConfigFormatNone ConfigFormat = iota
	ConfigFormatJSON
	ConfigFormatYAML
	ConfigFormatINI
	ConfigFormatEnvDir
)

// SupervisorMode* consts family defines possible work modes of the supervisor,
// supported by the guide-dog.
const (
	SupervisorModeNone   SupervisorMode = 0
	SupervisorModeSimple                = 1 << iota
	SupervisorModeRestarting
)

// Options is just a storage of the possible options with some interpretations.
type Options struct {
	ConfigFormat    ConfigFormat
	ConfigPath      string
	Debug           bool
	Envs            map[string]string
	GracefulTimeout time.Duration
	LockFile        *lockfile.Lock
	PathsToTrack    []string
	PTY             bool
	Signal          syscall.Signal
	Supervisor      SupervisorMode
}

func (opt *Options) String() string {
	return fmt.Sprintf("<Options(configFormat='%v', configPath='%v', pathsToTrack='%v', debug='%t', envs='%v', gracefulTimeout='%d', lockFile='%v', signal='%v', supervisor='%v')>",
		opt.ConfigFormat,
		opt.ConfigPath,
		opt.PathsToTrack,
		opt.Debug,
		opt.Envs,
		opt.GracefulTimeout,
		opt.LockFile,
		opt.Signal,
		opt.Supervisor)
}

func (sm SupervisorMode) String() string {
	mode := make([]string, 0, 2)

	if sm == SupervisorModeNone {
		mode = append(mode, "none")
	} else {
		if sm&SupervisorModeSimple > 0 {
			mode = append(mode, "simple")
		}
		if sm&SupervisorModeRestarting > 0 {
			mode = append(mode, "restarting")
		}
	}

	return strings.Join(mode, " / ")
}

func (cf ConfigFormat) String() string {
	switch cf {
	case ConfigFormatNone:
		return "none"
	case ConfigFormatJSON:
		return "json"
	case ConfigFormatYAML:
		return "yaml"
	case ConfigFormatINI:
		return "ini"
	case ConfigFormatEnvDir:
		return "envdir"
	default:
		return "ERROR"
	}
}

// NewOptions builds new Options struct based on the given parameter list
func NewOptions(debug bool,
	signal string,
	envs []string,
	gracefulTimeout time.Duration,
	configFormat string,
	configPath string,
	pathsToTrack []string,
	lockFile string,
	pty bool,
	supervise bool,
	restartOnConfigChanges bool) (options *Options, err error) {
	convertedConfigFormat, err := parseConfigFormat(configFormat)
	if err != nil {
		log.WithFields(log.Fields{
			"configFormat": configFormat,
			"error":        err,
		}).Errorf("Cannot convert configFormat.")
		return
	}

	convertedSignal, err := parseSignalName(signal)
	if err != nil {
		log.WithFields(log.Fields{
			"signal": signal,
			"error":  err,
		}).Errorf("Cannot convert signal.")
		return
	}

	convertedEnvs := parseEnvs(envs)

	supervisorMode := SupervisorModeNone
	if supervise {
		supervisorMode |= SupervisorModeSimple
	}
	if restartOnConfigChanges {
		supervisorMode |= SupervisorModeRestarting
	}

	var convertedLockFile *lockfile.Lock
	if lockFile != "" {
		convertedLockFile = lockfile.NewLock(lockFile)
	}

	options = &Options{
		ConfigFormat:    convertedConfigFormat,
		ConfigPath:      configPath,
		Debug:           debug,
		Envs:            convertedEnvs,
		GracefulTimeout: gracefulTimeout,
		LockFile:        convertedLockFile,
		PathsToTrack:    pathsToTrack,
		PTY:             pty,
		Signal:          convertedSignal,
		Supervisor:      supervisorMode,
	}

	return
}

func parseConfigFormat(name string) (format ConfigFormat, err error) {
	switch strings.ToLower(name) {
	case "":
		format = ConfigFormatNone
	case "json":
		format = ConfigFormatJSON
	case "yaml":
		format = ConfigFormatYAML
	case "ini":
		format = ConfigFormatINI
	case "envdir":
		format = ConfigFormatEnvDir
	default:
		err = fmt.Errorf("Unknown config format %s", name)
	}

	return
}

func parseSignalName(name string) (signal syscall.Signal, err error) {
	name = strings.ToUpper(name)
	if !strings.HasPrefix(name, "SIG") {
		name = "SIG" + name
	}

	switch name {
	case "SIGABRT":
		signal = syscall.SIGABRT
	case "SIGALRM":
		signal = syscall.SIGALRM
	case "SIGBUS":
		signal = syscall.SIGBUS
	case "SIGCHLD":
		signal = syscall.SIGCHLD
	case "SIGCLD":
		signal = syscall.SIGCLD
	case "SIGCONT":
		signal = syscall.SIGCONT
	case "SIGFPE":
		signal = syscall.SIGFPE
	case "SIGHUP":
		signal = syscall.SIGHUP
	case "SIGILL":
		signal = syscall.SIGILL
	case "SIGINT":
		signal = syscall.SIGINT
	case "SIGIO":
		signal = syscall.SIGIO
	case "SIGIOT":
		signal = syscall.SIGIOT
	case "SIGKILL":
		signal = syscall.SIGKILL
	case "SIGPIPE":
		signal = syscall.SIGPIPE
	case "SIGPOLL":
		signal = syscall.SIGPOLL
	case "SIGPROF":
		signal = syscall.SIGPROF
	case "SIGPWR":
		signal = syscall.SIGPWR
	case "SIGQUIT":
		signal = syscall.SIGQUIT
	case "SIGSEGV":
		signal = syscall.SIGSEGV
	case "SIGSTKFLT":
		signal = syscall.SIGSTKFLT
	case "SIGSTOP":
		signal = syscall.SIGSTOP
	case "SIGSYS":
		signal = syscall.SIGSYS
	case "SIGTERM":
		signal = syscall.SIGTERM
	case "SIGTRAP":
		signal = syscall.SIGTRAP
	case "SIGTSTP":
		signal = syscall.SIGTSTP
	case "SIGTTIN":
		signal = syscall.SIGTTIN
	case "SIGTTOU":
		signal = syscall.SIGTTOU
	case "SIGUNUSED":
		signal = syscall.SIGUNUSED
	case "SIGURG":
		signal = syscall.SIGURG
	case "SIGUSR1":
		signal = syscall.SIGUSR1
	case "SIGUSR2":
		signal = syscall.SIGUSR2
	case "SIGVTALRM":
		signal = syscall.SIGVTALRM
	case "SIGWINCH":
		signal = syscall.SIGWINCH
	case "SIGXCPU":
		signal = syscall.SIGXCPU
	case "SIGXFSZ":
		signal = syscall.SIGXFSZ
	default:
		err = fmt.Errorf("Unknown signal definition %s", name)
	}

	return
}

func parseEnvs(envs []string) (converted map[string]string) {
	converted = make(map[string]string)

	for _, env := range envs {
		split := strings.SplitN(env, "=", 2)
		if len(split) == 2 {
			converted[split[0]] = split[1]
		} else {
			converted[split[0]] = ""
		}
	}

	return
}
