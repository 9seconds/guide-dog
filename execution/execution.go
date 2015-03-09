package execution

import (
	environment "github.com/9seconds/guide-dog/environment"
)

func Execute(command []string, env *environment.Environment) int {
	pathsToWatch := []string{env.Options.ConfigPath}
	for _, path := range env.Options.PathsToTrack {
		pathsToWatch = append(pathsToWatch, path)
	}

	watcherChannel := makeWatcher(pathsToWatch)
	defer close(watcherChannel)

	exitCodeChannel := make(chan int, 1)
	defer close(exitCodeChannel)

	return <-exitCodeChannel
}
