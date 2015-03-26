package options

import (
	"fmt"
	"strings"
	"syscall"
)

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
	case "SIGPROF":
		signal = syscall.SIGPROF
	case "SIGQUIT":
		signal = syscall.SIGQUIT
	case "SIGSEGV":
		signal = syscall.SIGSEGV
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
