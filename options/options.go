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
	ConfigFormat   uint8
	SupervisorMode uint8
)

const (
	CONFIG_FORMAT_NONE ConfigFormat = iota
	CONFIG_FORMAT_JSON
	CONFIG_FORMAT_YAML
	CONFIG_FORMAT_INI
	CONFIG_FORMAT_ENVDIR
)

const (
	SUPERVISOR_MODE_NONE SupervisorMode = 1 << iota
	SUPERVISOR_MODE_SIMPLE
	SUPERVISOR_MODE_RESTARTING
)

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

	if sm == SUPERVISOR_MODE_NONE {
		mode = append(mode, "none")
	} else {
		if sm&SUPERVISOR_MODE_SIMPLE > 0 {
			mode = append(mode, "simple")
		}
		if sm&SUPERVISOR_MODE_RESTARTING > 0 {
			mode = append(mode, "restarting")
		}
	}

	return strings.Join(mode, " / ")
}

func (cf ConfigFormat) String() string {
	switch cf {
	case CONFIG_FORMAT_NONE:
		return "none"
	case CONFIG_FORMAT_JSON:
		return "json"
	case CONFIG_FORMAT_YAML:
		return "yaml"
	case CONFIG_FORMAT_INI:
		return "ini"
	case CONFIG_FORMAT_ENVDIR:
		return "envdir"
	default:
		return "ERROR"
	}
}

func NewOptions(debug bool, signal string, envs []string,
	gracefulTimeout time.Duration, configFormat string,
	configPath string, pathsToTrack []string, lockFile string,
	pty bool, supervise bool, restartOnConfigChanges bool) (options *Options, err error) {
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

	supervisorMode := SUPERVISOR_MODE_NONE
	if supervise {
		supervisorMode |= SUPERVISOR_MODE_SIMPLE
	}
	if restartOnConfigChanges {
		supervisorMode |= SUPERVISOR_MODE_RESTARTING
	}

	var convertedLockFile *lockfile.Lock = nil
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
		format = CONFIG_FORMAT_NONE
	case "json":
		format = CONFIG_FORMAT_JSON
	case "yaml":
		format = CONFIG_FORMAT_YAML
	case "ini":
		format = CONFIG_FORMAT_INI
	case "envdir":
		format = CONFIG_FORMAT_ENVDIR
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
