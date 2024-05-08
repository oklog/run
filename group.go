package run

import (
	"context"
	"time"
)

// group collects actors (functions) and runs them concurrently.
// When one actor (function) returns, all actors are interrupted.
// The zero value of a Group is useful.
type group struct {
	actors       []actor
	closeTimeout time.Duration
	syncShutdown bool
}

// Add an actor (function) to the group. Each actor must be pre-emptable by an
// interrupt function. That is, if interrupt is invoked, execute should return.
// Also, it must be safe to call interrupt even after execute has returned.
//
// The first actor (function) to return interrupts all running actors.
// The error is passed to the interrupt functions, and is returned by Run.
func (g *group) add(execute func(context.Context) error, interrupt func(context.Context) error) {
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

	runCtx, runCancel := context.WithCancel(context.Background())

	// Run each actor.
	runCh := make(chan error, len(g.actors))
	defer close(runCh)

	for _, a := range g.actors {
		go func(a actor) {
			runCh <- a.execute(runCtx)
		}(a)
	}

	// Wait for the first actor to stop.
	err := <-runCh

	// Notify Run() that is needs to stop.
	runCancel()

	var closeCtx context.Context
	{
		if g.closeTimeout == 0 {
			closeCtx = context.Background()
		} else {
			ctx, cancel := context.WithTimeout(context.Background(), g.closeTimeout)
			defer cancel()
			closeCtx = ctx
		}
	}

	// Notify Close() that it needs to stop.
	closeCh := make(chan struct{}, len(g.actors))
	defer close(closeCh)

	for _, a := range g.actors {
		a := a // NOTE(frank): May not need this anymore in go1.22.

		shutdown := func(a actor) {
			a.interrupt(closeCtx)
			closeCh <- struct{}{}
		}

		if g.syncShutdown {
			shutdown(a)
		} else {
			go shutdown(a)
		}
	}

	// Wait for all Close() to stop.
	for i := 0; i < cap(closeCh); i++ {
		<-closeCh
	}

	// Wait for all actors to stop.
	for i := 1; i < cap(runCh); i++ {
		<-runCh
	}

	// Return the original error.
	return err
}

type actor struct {
	execute   func(context.Context) error
	interrupt func(context.Context) error
}
