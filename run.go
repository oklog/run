package run

import (
	"go.uber.org/zap"
)

const (
	Always = true
	Never  = false
)

type Group struct {
	runnables []Runnable
	group     group
	options   *Options
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
	Run() error

	// Close is responsible for gracefully shutting down the component. It takes an error
	// as an input parameter, which can be used to provide information about the reason for
	// the shutdown. The Close method should ensure that all resources used by the component
	// are properly released and any necessary cleanup is performed. If the component is
	// successfully closed, it should return nil. If an error occurs during the shutdown
	// process, it should return the error.
	//
	// This method is REQUIRED.
	Close(error) error

	// Alive assesses whether the Runnable has properly been initialized and is ready to perform.
	// It does not assess health.
	//
	// This method is REQUIRED.
	Alive() bool

	// Name returns the name of the Runnable.
	//
	// This method is OPTIONAL.
	Name() string

	// Fields returns the log fields of the Runnable.
	//
	// This method is OPTIONAL.
	Fields() []zap.Field

	// compatibility forces clients to include ForwardCompatibility in their Runnable implementations.
	// This is so we can make forward compatible changes to the Runnable interface.
	compatibility()
}

func New(options ...Option) *Group {
	defaults := &Options{
		logger: zap.NewNop(),
	}

	for _, fn := range options {
		fn(defaults)
	}

	return &Group{
		options: defaults,
	}
}

func (g *Group) Add(condition bool, r Runnable) {
	if !condition {
		return
	}

	logger := g.options.logger.With(append([]zap.Field{zap.String("name", r.Name())}, r.Fields()...)...)
	g.runnables = append(g.runnables, r)

	g.group.Add(func() error {
		logger.Info("starting runnable")
		defer func() {
			logger.Info("stopped runnable")
		}()

		return r.Run()
	}, func(err error) {
		logger.Info("closing runnable")
		defer func() {
			logger.Info("stopping runnable")
		}()

		r.Close(err)
	})
}

func (g *Group) Run() error {
	return g.group.Run()
}

func (g *Group) Alive() bool {
	for _, r := range g.runnables {
		if !r.Alive() {
			return false
		}
	}

	return true
}
