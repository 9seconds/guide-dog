package execution

import (
	log "github.com/Sirupsen/logrus"
	fsnotify "gopkg.in/fsnotify.v1"
)

func makeWatcher(paths []string) (channel chan bool) {
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

	go func() {
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
