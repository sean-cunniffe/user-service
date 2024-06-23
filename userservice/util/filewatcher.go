package util

// a package for watching directories for changes using fsnotify
import (
	"context"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

// watch starts watching the directories for changes
func WatchDir(ctx context.Context, onChange func(string), dirs ...string) error {
	watcher, err := fsnotify.NewWatcher()
	for _, dir := range dirs {
		err := watcher.Add(dir)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-watcher.Events:
				if !ok {
					log.Error("error reading from file watcher events")
				}
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					onChange(event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					log.Errorf("error reading from file watcher errors: %v", err)
				}
			}
		}
	}()
	return nil
}
