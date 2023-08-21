package tidy

import (
	"github.com/duexcoast/tidy-up/pkg/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
)

// TODO: I'm thinking this function should be conditionally called in the Tidy.Sort method, only
// if the watch flag is set. We can then pass in a callback function to use the correct sorting
// method within this function.
func fileWatcher(sortDir string, sortFunc func() error) error {
	logger := logger.Get()
	// create a new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	// start listening for events
	go func() {
		// TODO: refactor to use range instead of for select loop
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				logger.Info().Str("Fs event:", event.Name)
				if event.Has(fsnotify.Create) {
					// a file has been created: we want to sort the sortDir here.
					// We do so using the passed in callback function - sortFunc()
					logger.Info().Str("file created: ", event.Name)
					sortFunc()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Error().AnErr("Error watching file:", err)
			}
		}
	}()

	// add a path
	err = watcher.Add(sortDir)
	if err != nil {
		return err
	}

	// TODO: ordinarily we would block forever if we were running this in a main go routine.
	// Where should we block? Should it still be here or do we want to return from this code.
	// because we have more to run after it.
	<-make(chan struct{})
	return nil
}
