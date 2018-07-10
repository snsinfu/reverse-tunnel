package taskch

import (
	"sync"
)

// T is the taskch type that manages collection of concurrent tasks. A task is
// a function, executed as a goroutine, returning an error to indicate success
// or failure.
type T struct {
	wg   *sync.WaitGroup
	join chan bool
	err  chan error
}

// New creates a taskch object.
func New() *T {
	tch := T{
		wg:   &sync.WaitGroup{},
		join: make(chan bool),
		err:  make(chan error),
	}

	go func() {
		tch.wg.Wait()
		tch.join <- true
	}()

	return &tch
}

// Go launches a task as a goroutine.
func (tch *T) Go(task func() error) {
	tch.wg.Add(1)

	go func() {
		defer tch.wg.Done()

		if err := task(); err != nil {
			tch.err <- err
		}
	}()
}

// Wait blocks until all the tasks to success or any of the tasks to fail. In
// the latter case, one of the errors is returned. To ensure that all tasks to
// finish, repeatedly call Wait until nil is returned.
func (tch *T) Wait() error {
	select {
	case err := <-tch.err:
		return err
	case <-tch.join:
		return nil
	}
}
