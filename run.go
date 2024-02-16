package run

import (
	"context"
	"log/slog"
	"os"
)

type Group struct {
	runnables []Runnable
	group     group
	logger    *slog.Logger
}

// Runnable greatly simplifies the process of propagating graceful shutdown signals
// between multiple components within a Go program. By implementing the Runnable
// interface for each component, it becomes easy to manage and orchestrate their
// execution and termination in a clean and deterministic manner. When a Runnable
// component (e.g., Runnable A) needs to exit, the Close method can be called with
// an appropriate error or reason for the shutdown. This makes it straightforward to
// inform all other Runnable components that they also need to shut down gracefully.
// The caller can then iterate through all the components, calling their respective
// Close methods, and ensuring that each component is aware of the shutdown signal
// and can perform the necessary cleanup and resource release.
//
// This approach promotes a clean and organized shutdown process, allowing for better
// resource management and error handling. By using the Runnable interface, developers
// can create more robust and maintainable Go programs, with clear and deterministic
// control over the lifecycle of each component.
type Runnable interface {
	// Run is responsible for starting the execution of the component. It should contain
	// the main logic of the component and is expected to run until an error occurs or
	// the component is stopped. If the component runs successfully, it should return nil.
	// If an error occurs during the execution, it should return the error.
	//
	// This method is REQUIRED.
	Run(context.Context) error

	// Close is responsible for gracefully shutting down the component. It takes an error
	// as an input parameter, which can be used to provide information about the reason for
	// the shutdown. The Close method should ensure that all resources used by the component
	// are properly released and any necessary cleanup is performed. If the component is
	// successfully closed, it should return nil. If an error occurs during the shutdown
	// process, it should return the error.
	//
	// This method is OPTIONAL.
	Close()

	// Alive assesses whether the Runnable has properly been initialized and is ready to perform.
	// It does not assess health.
	//
	// This method is OPTIONAL.
	Alive() bool

	// Name returns the name of the Runnable.
	//
	// This method is OPTIONAL.
	Name() string

	// Fields returns the log fields of the Runnable.
	//
	// This method is OPTIONAL.
	Fields() []slog.Attr
}

// New is syntactic sugar for creating a new
// Group with the provided functional options.
func New(options ...Option) *Group {
	defaults := []Option{
		WithLogger(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})),
	}

	g := &Group{}
	for _, fn := range append(defaults, options...) {
		fn(g)
	}

	return g
}

// Add appends each Runnable to the Group if the condition is met.
func (g *Group) Add(when bool, runnables ...Runnable) {
	if !when {
		return
	}

	for _, r := range runnables {
		g.add(r)
	}
}

func (g *Group) add(r Runnable) {
	logger := g.logger.With(fields(append([]slog.Attr{slog.String("name", r.Name())}, r.Fields()...))...)
	g.runnables = append(g.runnables, r)

	g.group.add(func(ctx context.Context) error {
		logger.Info("started", "method", "run")
		defer func() {
			logger.Info("returned", "method", "run")
		}()

		return r.Run(ctx)
	}, func() {
		logger.Info("started", "method", "close")
		defer func() {
			logger.Info("returned", "method", "close")
		}()

		r.Close()
	})
}

// Run all Runnables concurrently. When the first Runnable returns,
// all others are interrupted. Run returns when all Runnables have
// exited. The first encountered error is returned.
func (g *Group) Run() error {
	return g.group.run()
}

// Alive asseses the liveness of all registered Runnables.
func (g *Group) Alive() bool {
	for _, r := range g.runnables {
		if !r.Alive() {
			return false
		}
	}

	return true
}

func fields(attr []slog.Attr) []any {
	out := make([]any, 0, len(attr)*2)

	for _, a := range attr {
		out = append(out, a.Key, a.Value.Any())
	}

	return out
}
