package build

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// App the current app we want to build (never changes so we can make it static)
var App = &Target{}

func init() {
	var err error
	bind := randomLocalBind()

	// We want to build on the first request
	App.rwmut = &sync.RWMutex{}
	App.rebuild = true

	App.Bind = bind
	App.URL, err = url.Parse(fmt.Sprintf("http://%s", App.Bind))
	if err != nil { // This should not happen
		log.Fatal(err)
	}

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	binName := filepath.Base(path)

	if runtime.GOOS == "windows" {
		App.BinaryPath = filepath.Join(os.TempDir(), binName+".exe")
	} else {
		App.BinaryPath = filepath.Join(os.TempDir(), binName)
	}

	// Monitor the filesytem an flag for rebuild on a change
	go func() {
		for _ = range FileChanged("./") {
			App.Rebuild()
		}
	}()
}

func randomLocalBind() string {
	max := 65536
	min := 9000
	random := rand.New(rand.NewSource(time.Now().Unix()))

	return fmt.Sprintf("localhost:%d", random.Intn(max-min)+min)
}
