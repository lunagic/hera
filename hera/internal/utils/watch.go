package utils

import (
	"log"
	"time"

	"github.com/radovskyb/watcher"
)

func Watch(worker func(name string) error) {

	w := watcher.New()

	w.SetMaxEvents(1)

	go func() {
		for {
			select {
			case event := <-w.Event:
				_ = worker(event.Path)
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.AddRecursive("."); err != nil {
		log.Fatalln(err)
	}

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}
