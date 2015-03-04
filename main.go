package main

import (
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v1"
)

var (
	cmdLine = kingpin.New("guide-dog", "Small supervisor with envdir possibilities")

	debug = cmdLine.Flag("debug", "Enable debug mode").Bool()

	signal       = cmdLine.Flag("signal", "Signal to graceful timeout given process").Default("TERM").String()
	configFormat = cmdLine.Flag("config-format", "Format of configs").String()
	configFile   = cmdLine.Flag("config-path", "Config path").String()
	lockFile     = cmdLine.Flag("lock-file", "Lockfile on the local machine to acquire").String()
)

func main() {
	kingpin.MustParse(cmdLine.Parse(os.Args[1:]))
}
