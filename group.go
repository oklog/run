package run

import (
	"context"
)

// group collects actors (functions) and runs them concurrently.
// When one actor (function) returns, all actors are interrupted.
// The zero value of a Group is useful.
type group struct {
	actors []actor
}

// Add an actor (function) to the group. Each actor must be pre-emptable by an
// interrupt function. That is, if interrupt is invoked, execute should return.
// Also, it must be safe to call interrupt even after execute has returned.
//
// The first actor (function) to return interrupts all running actors.
// The error is passed to the interrupt functions, and is returned by Run.
func (g *group) add(execute func(context.Context) error, interrupt func()) {
	g.actors = append(g.actors, actor{execute, interrupt})
}

// Run all actors (functions) concurrently.
// When the first actor returns, all others are interrupted.
// Run only returns when all actors have exited.
// Run returns the error returned by the first exiting actor.
func (g *group) run() error {
	if len(g.actors) == 0 {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Run each actor.
	run := make(chan error, len(g.actors))
	for _, a := range g.actors {
		go func(a actor) {
			run <- a.execute(ctx)
		}(a)
	}

	// Wait for the first actor to stop.
	err := <-run

	// Notify Run() that is needs to stop.
	cancel()

	// Notify Close() that it needs to stop.
	close := make(chan struct{})
	for _, a := range g.actors {
		go func(a actor) {
			a.interrupt()
			close <- struct{}{}
		}(a)
	}

	// Wait for all Close() to stop.
	for i := 1; i < cap(close); i++ {
		<-close
	}

	// Wait for all actors to stop.
	for i := 1; i < cap(run); i++ {
		<-run
	}

	// Return the original error.
	return err
}

type actor struct {
	execute   func(context.Context) error
	interrupt func()
}
