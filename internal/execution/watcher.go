// Package execution contains all logic for execution of external commands
// based on Environment struct.
//
// This file contains definition of filesystem notifications tracking.
package execution

import (
	log "github.com/Sirupsen/logrus"
	fsnotify "gopkg.in/fsnotify.v1"

	environment "github.com/9seconds/guide-dog/internal/environment"
)

// makeWatcher starts to track given paths and sends filesystem notifications
// into channel.
func makeWatcher(paths []string, env *environment.Environment) (channel chan bool) {
	channel = make(chan bool, 1)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	trackPaths := make([]string, 0, len(paths))
	for _, value := range paths {
		if value != "" {
			trackPaths = append(trackPaths, value)
		}
	}

	if len(trackPaths) == 0 {
		return
	}

	for _, path := range trackPaths {
		log.WithField("path", path).Info("Add path")
		err = watcher.Add(path)
		if err != nil {
			log.WithFields(log.Fields{
				"path":  path,
				"error": err,
			}).Warn("Cannot add path")
		}
	}

	go watcherLoop(env, channel, watcher)

	return
}

// watcherLoop defines main watcher loop.
func watcherLoop(env *environment.Environment, channel chan bool, watcher *fsnotify.Watcher) {
	defer watcher.Close()

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Op == 0 {
				continue
			}

			log.WithFields(log.Fields{
				"event": event,
				"op":    event.Op,
			}).Info("Event from filesystem is coming")

			env.Update()

			if len(channel) == 0 {
				channel <- true
			}
		case err := <-watcher.Errors:
			if err != nil {
				log.WithField("error", err).Error("Some problem with filesystem notifications")
			}
		}
	}
}
