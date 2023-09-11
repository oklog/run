package run

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

// ContextHandler returns an actor, i.e. an execute and interrupt func, that
// terminates with when the parent context is canceled.
func ContextHandler(ctx context.Context, signals ...os.Signal) (execute func() error, interrupt func(error)) {
	ctx, cancel := context.WithCancel(ctx)
	return func() error {
			<-ctx.Done()
			return ctx.Err()
		}, func(error) {
			cancel()
		}
}

// SignalHandler returns an actor, i.e. an execute and interrupt func, that
// terminates with SignalError when the process receives one of the provided
// signals, or the parent context is canceled. If no signals are provided,
// handler will listen for all incoming signals.
func SignalHandler(ctx context.Context, signals ...os.Signal) (execute func() error, interrupt func(error)) {
	ctx, cancel := context.WithCancel(ctx)
	return func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, signals...)
			defer signal.Stop(c)
			select {
			case sig := <-c:
				return SignalError{Signal: sig}
			case <-ctx.Done():
				return ctx.Err()
			}
		}, func(error) {
			cancel()
		}
}

// SignalError is returned by the signal handler's execute function
// when it terminates due to a received signal.
type SignalError struct {
	Signal os.Signal
}

// Error implements the error interface.
func (e SignalError) Error() string {
	return fmt.Sprintf("received signal %s", e.Signal)
}
