package utils

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/radovskyb/watcher"
)

func Watch(paths []string, ignores []string, worker func(name string) error) {
	ignores = append(ignores, "/.git")

	w := watcher.New()

	w.SetMaxEvents(1)

	w.AddFilterHook(func(info os.FileInfo, fullPath string) error {
		for _, ignore := range ignores {
			if strings.Contains(fullPath, ignore) {
				return watcher.ErrSkip
			}
		}

		return nil
	})

	go func() {
		for {
			select {
			case event := <-w.Event:
				worker(event.Path)
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	for _, path := range paths {
		if err := w.AddRecursive(path); err != nil {
			log.Fatalln(err)
		}

	}

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}
