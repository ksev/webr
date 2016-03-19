package build

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/ksev/webr/config"
)

// Target Information about a Go web app to build and monitor
type Target struct {
	URL        *url.URL
	Bind       string
	BinaryPath string
	rebuild    bool
	err        error
	rwmut      *sync.RWMutex
	command    *exec.Cmd
}

// Rebuild signal to the target that it need to rebuild
func (t *Target) Rebuild() {
	t.rwmut.Lock()
	defer t.rwmut.Unlock()

	t.rebuild = true
}

// Build build the target, blocks while building
func (t *Target) Build() error {
	build := exec.Command("go", "build", "-o", t.BinaryPath)
	data, err := build.CombinedOutput()

	if err != nil {
		return errors.New(string(data))
	}

	return nil
}

// Reachable see if the underlying web app is reachable with http
func (t *Target) Reachable() bool {
	resp, err := http.Get(t.URL.String())
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return true
}

// Run run the target, will error out it has not yet been built.
// blocks until its reachable via the Target bind or target quits
func (t *Target) Run() error {
	args := append(config.Current.TargetPassTrough, config.Current.TargetBindArgName, t.Bind)
	run := exec.Command(t.BinaryPath, args...)

	run.Stdout = os.Stdout
	run.Stderr = os.Stderr

	if err := run.Start(); err != nil {
		return err
	}

	unblock := make(chan error, 5)
	stop := make(chan bool, 5)

	go func() {
	Loop:
		for {
			select {
			case <-stop:
				break Loop
			default:
				if t.Reachable() {
					unblock <- nil
					stop <- true
					continue Loop
				}
				time.Sleep(5 * time.Millisecond)
			}
		}
	}()

	go func() {
		err := run.Wait()
		unblock <- err
		stop <- true
	}()

	t.command = run

	return <-unblock
}

// CheckAndWait Check if the Target need to be rebuild and restarted, block while waiting
func (t *Target) CheckAndWait() error {
	t.rwmut.RLock()
	if !t.rebuild {
		t.rwmut.RUnlock()
		return t.err
	}
	t.rwmut.RUnlock()

	t.rwmut.Lock()
	defer t.rwmut.Unlock()
	if !t.rebuild {
		return t.err
	}
	t.rebuild = false

	if t.command != nil {
		t.command.Process.Kill()
	}

	if err := t.Build(); err != nil {
		t.err = err
		return err
	}

	if err := t.Run(); err != nil {
		t.err = err
		return err
	}

	t.err = nil
	return nil
}
