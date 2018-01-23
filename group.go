// Package run implements an actor-runner with deterministic teardown. It is
// somewhat similar to package errgroup, except it does not require actor
// goroutines to understand context semantics. This makes it suitable for use in
// more circumstances; for example, goroutines which are handling connections
// from net.Listeners, or scanning input from a closable io.Reader.
package run

// Group collects actors (functions) and runs them concurrently.
// When one actor (function) returns, except for a sidecar actor
// that does not return an error, all actors are interrupted.
// The zero value of a Group is useful.
type Group struct {
	actors []actor
}

// Add an actor (function) to the group. Each actor must be pre-emptable by an
// interrupt function. That is, if interrupt is invoked, execute should return.
// Also, it must be safe to call interrupt even after execute has returned.
//
// The first actor (function) to return interrupts all running actors.
// The error is passed to the interrupt functions, and is returned by Run.
func (g *Group) Add(execute func() error, interrupt func(error)) {
	g.actors = append(g.actors, actor{
		execute:   execute,
		interrupt: interrupt,
		sidecar:   false,
	})
}

// AddSidecar add a sidecar actor to the group. Each actor must satisfy
// the same conditions as a normal actor from Add function.
//
// If the sidecar actor does not return an error, it does not interrupt
// other running actors. Otherwise it acts like a normal actor.
func (g *Group) AddSidecar(execute func() error, interrupt func(error)) {
	g.actors = append(g.actors, actor{
		execute:   execute,
		interrupt: interrupt,
		sidecar:   true,
	})
}

// Run all actors (functions) concurrently.
// There are two cases when all actors will be interrupted.
// The first case, when the first actor returns an error.
// The second case, when the first normal actor returns.
// Run only returns when all actors have exited.
// Run returns the error returned by the first exiting actor.
func (g *Group) Run() error {
	if len(g.actors) == 0 {
		return nil
	}

	// Run each actor.
	actorExecutionResultChan := make(chan *actorExecutionResult, len(g.actors))
	for _, a := range g.actors {
		go func(a actor) {
			actorExecutionResultChan <- &actorExecutionResult{
				fromSidecar: a.sidecar,
				err:         a.execute(),
			}
		}(a)
	}

	// The original error.
	var err error

	// Number of actors that need to wait.
	waitCount := len(g.actors)

	// Wait for the first normal actor to stop
	// or the sidecar actor returns an error.
	for waitCount > 0 {
		execResult := <-actorExecutionResultChan
		if execResult.fromSidecar && execResult.err == nil {
			// If sidecar actor stopped without error, then
			// we do not need to wait for it to stop.
			waitCount--
		} else {
			err = execResult.err
			break
		}
	}

	// Signal all actors to stop.
	for _, a := range g.actors {
		a.interrupt(err)
	}

	// Wait for actors to stop.
	for i := 1; i < waitCount; i++ {
		<-actorExecutionResultChan
	}

	// Return the original error.
	return err
}

type actor struct {
	execute   func() error
	interrupt func(error)
	sidecar   bool
}

type actorExecutionResult struct {
	fromSidecar bool
	err         error
}
