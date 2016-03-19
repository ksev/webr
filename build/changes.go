package build

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// FileChanged get notified when a file has changed in any subfolder
func FileChanged(path string) <-chan string {
	modified := make(chan string, 500)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
	Loop:
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					continue Loop
				}

				switch event.Op {
				case event.Op & fsnotify.Create:
					fi, err := os.Stat(event.Name)
					if err != nil {
						log.Panic(err)
					}
					if fi.IsDir() {
						watcher.Add(event.Name)
					}
				case event.Op & fsnotify.Remove:
					watcher.Remove(event.Name)
				}

				modified <- event.Name

			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	// Add all folders recursivly
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			err = watcher.Add(path)
			if err != nil {
				log.Fatal(err)
			}
		}
		return err
	})

	return modified
}
