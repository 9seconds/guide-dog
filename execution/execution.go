package execution

import (
	environment "github.com/9seconds/guide-dog/environment"
)

func Execute(command []string, env *environment.Environment) int {
	exitCodeChannel := make(chan int, 1)

	watcher, watcherChannel := makeWatcher(env.Options.ConfigPath)
	defer close(watcherChannel)
	defer watcher.Close()

	return <-exitCodeChannel
}
