// Package run does something.
package run

import (
	"context"
	"log/slog"
	"os"
	"time"
)

var (
	defaultLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
)

// Group manages a collection of Runnables.
type Group struct {
	runnables    []Runnable
	group        group
	logger       *slog.Logger
	closeTimeout time.Duration
}

type Runnable interface {
	// Run is responsible for executing the main logic of the component
	// and is expected to run until it needs to shut down. Cancellation
	// of the provided context signals that the component should shut down
	// gracefully. If the Runnable implements the Close method than this
	// context can be ignored.
	//
	// Implementations must insure that instantiations of things to be
	// shutdown do not leak outside of this method (i.e. a constructor
	// calling net.Listen) as the Close method may not be called.
	Run(context.Context) error

	// Close is responsible for gracefully shutting down the component. It
	// can either initiate the shutdown process and return or wait for the
	// shutdown process to complete before returning. This method should
	// ensure that all resources used by the component are properly released
	// and any necessary cleanup is performed. If this method is not implemented
	// it is expected that the Run method properly handle context cancellation.
	Close(context.Context) error

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
		WithLogger(defaultLogger),
	}

	g := &Group{}
	for _, fn := range append(defaults, options...) {
		fn(g)
	}

	g.group.closeTimeout = g.closeTimeout

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
	// In case clients do not use the New function to create a new Group.
	if g.logger == nil {
		g.logger = defaultLogger
	}

	logger := g.logger.With(
		anything(append([]slog.Attr{
			slog.String("name", r.Name()),
		}, r.Fields()...))...,
	)

	g.runnables = append(g.runnables, r)

	g.group.add(func(ctx context.Context) error {
		return do(ctx, r.Run, logger.With(slog.String("method", "run")))
	}, func(ctx context.Context) error {
		return do(ctx, r.Close, logger.With(slog.String("method", "close")))
	})
}

func do(ctx context.Context, fn func(context.Context) error, logger *slog.Logger) error {
	logger.Info("started")

	if err := fn(ctx); err != nil {
		logger.Error("failed", slog.String("error", err.Error()))
		return err
	}

	logger.Info("returned")

	return nil
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
