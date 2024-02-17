// Package run does something.
package run

import (
	"context"
	"log/slog"
	"os"
)

// Group manages a collection of Runnables.
type Group struct {
	runnables []Runnable
	group     group
	logger    *slog.Logger
}

type Runnable interface {
	// Run is responsible for executing the main logic of the component
	// and is expected to run until it needs to shut down. Cancellation
	// of the provided context signals that the component should shut down
	// gracefully. If the Runnable implements the Close method than this
	// context can be ignored.
	Run(context.Context) error

	// Close is responsible for gracefully shutting down the component. It
	// can either initiate the shutdown process and return or wait for the
	// shutdown process to complete before returning. This method should
	// ensure that all resources used by the component are properly released
	// and any necessary cleanup is performed. If this method is not implemented
	// it is expected that the Run method properly handle context cancellation.
	Close()

	// Alive assesses whether the Runnable has
	// properly been initialized and is active.
	Alive() bool

	// Name returns the name of the Runnable.
	Name() string

	// Fields allows clients to attach additional fields
	// to every log message this library produces.
	Fields() []slog.Attr
}

// New is syntactic sugar for creating a new
// Group with the provided functional options.
func New(options ...Option) *Group {
	defaults := []Option{
		WithLogger(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					return slog.Attr{}
				}

				return a
			},
		}))),
	}

	g := &Group{}
	for _, fn := range append(defaults, options...) {
		fn(g)
	}

	return g
}

// Add appends each Runnable to the group if the condition is met.
func (g *Group) Add(when bool, runnables ...Runnable) {
	if !when {
		return
	}

	for _, r := range runnables {
		g.add(r)
	}
}

// Always adds each Runnable to the group.
func (g *Group) Always(runnables ...Runnable) {
	g.Add(true, runnables...)
}

func (g *Group) add(r Runnable) {
	logger := g.logger.With(
		anything(append([]slog.Attr{
			slog.String("name", r.Name()),
		}, r.Fields()...))...,
	)

	g.runnables = append(g.runnables, r)

	g.group.add(func(ctx context.Context) error {
		logger.Info("started", slog.String("method", "run"))
		defer func() {
			logger.Info("returned", slog.String("method", "run"))
		}()

		return r.Run(ctx)
	}, func() {
		logger.Info("started", slog.String("method", "close"))
		defer func() {
			logger.Info("returned", slog.String("method", "close"))
		}()

		r.Close()
	})
}

// Run invokes and manages all registered Runnables.
//
//  1. Invoke Run on each Runnable concurrently.
//  2. Wait for the first Runnable to return.
//  3. Cancel the context passed to Run.
//  4. Invoke Close on each Runnable concurrently.
//  5. Wait for all Close methods to return.
//  6. Wait for all Run methods to return.
//
// It returns the initial error.
func (g *Group) Run() error {
	return g.group.run()
}

// Alive assess the liveness of all registered Runnables.
func (g *Group) Alive() bool {
	for _, r := range g.runnables {
		if !r.Alive() {
			return false
		}
	}

	return true
}

func anything[T any](data []T) []any {
	out := make([]any, len(data))

	for i, v := range data {
		out[i] = v
	}

	return out
}
