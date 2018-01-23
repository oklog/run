// Package run implements an actor-runner with deterministic teardown. It is
// somewhat similar to package errgroup, except it does not require actor
// goroutines to understand context semantics. This makes it suitable for use in
// more circumstances; for example, goroutines which are handling connections
// from net.Listeners, or scanning input from a closable io.Reader.
package run

// Group collects actors (functions) and runs them concurrently.
// There are two types of actors (functions): regular and sidecar.
// They differ in how their returns are handled. When a regular
// actor (function) returns, all actors are interrupted unconditionally.
// When a sidecar actor (function) exits, all actors are interrupted only
// if an error is returned. The zero value of a Group is useful.
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

// Add a sidecar actor to the group. Each actor must meet the same
// requirements as a regular actor as desribed in the documentation for
// Add function.
//
// If a sidecar actor does not return an error, it does not interrupt
// all running actors. Otherwise it acts like a regular actor.
func (g *Group) AddSidecar(execute func() error, interrupt func(error)) {
	g.actors = append(g.actors, actor{
		execute:   execute,
		interrupt: interrupt,
		sidecar:   true,
	})
}

// Run all actors (functions) concurrently. It interrupts all actors when
// either a regular actor returns, or a sidecar actor returns an error.
// Run only returns when all actors have exited. Run returns the error
// returned by the first exiting actor.
func (g *Group) Run() error {
	if len(g.actors) == 0 {
		return nil
	}

	// Run each actor.
	actorResults := make(chan actorResult, len(g.actors))
	for _, a := range g.actors {
		go func(a actor) {
			actorResults <- actorResult{
				sidecar: a.sidecar,
				err:     a.execute(),
			}
		}(a)
	}

	// An original error to return.
	var err error

	// A number of actors that need to be waited to return.
	waitCount := len(g.actors)

	// Wait till a) a regular actor finishes, or b) a sidecar actor
	// returns an error, or c) all sidecar actors returno with no error, if
	// there are no regular actors.
	for waitCount > 0 {
		result := <-actorResults
		if result.sidecar && result.err == nil {
			// If a sidecar actor stops without an error, then
			// we do not need to wait for it to stop later.
			waitCount--
		} else {
			// Remember the error to return from this function.
			err = result.err
			break
		}
	}

	// Signal all actors to stop.
	for _, a := range g.actors {
		a.interrupt(err)
	}

	// Wait for actors to stop.
	for i := 1; i < waitCount; i++ {
		<-actorResults
	}

	// Return the original error.
	return err
}

type actor struct {
	execute   func() error
	interrupt func(error)
	sidecar   bool
}

type actorResult struct {
	sidecar bool
	err     error
}
