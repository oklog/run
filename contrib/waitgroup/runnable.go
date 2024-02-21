package waitgroup

import (
	"context"
	"sync"

	"github.com/superblocksteam/run"
)

type wait struct {
	done chan struct{}
	wg   *sync.WaitGroup

	run.ForwardCompatibility
}

// NewWait returns a runnable that ensures that the wait group completes.
// This is useful when you want to wait for dynamically created tasks
// (i.e. async api executions) to complete before exiting.
func NewWait(wg *sync.WaitGroup) run.Runnable {
	return &wait{
		wg:   wg,
		done: make(chan struct{}),
	}
}

func (w *wait) Run(context.Context) error {
	<-w.done
	w.wg.Wait()
	return nil
}

func (w *wait) Name() string { return "wait group reaper" }

func (w *wait) Close(context.Context) error {
	close(w.done)
	return nil
}
