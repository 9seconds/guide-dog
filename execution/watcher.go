package execution

import (
	log "github.com/Sirupsen/logrus"
	fsnotify "gopkg.in/fsnotify.v1"
)

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
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				if err != nil {
					log.WithField("error", err).Error("Some problem with filesystem notifications")
				}
			}
		}
	}()

	return
}
