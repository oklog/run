package run

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

// SignalHandler returns an actor, i.e. an execute and interrupt func, that
// terminates with SignalError when the process receives one of the provided
// signals, or the parent context is canceled.
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

// Upgrader handles zero downtime upgrades and passing files between processes.
//
// Upgrader is based on https://github.com/cloudflare/tableflip.
type Upgrader interface {
	Ready() error
	Exit() <-chan struct{}
	Stop()
}

// GracefulRestart returns an actor, i.e. an execute and interrupt func, that
// terminates when graceful restart is initiated and the child process
// signals to be ready, or the parent context is canceled.
func GracefulRestart(ctx context.Context, upg Upgrader) (execute func() error, interrupt func(error)) {
	ctx, cancel := context.WithCancel(ctx)

	return func() error {
			// Tell the parent we are ready
			err := upg.Ready()
			if err != nil {
				return err
			}

			select {
			case <-upg.Exit(): // Wait for child to be ready (or application shutdown)
				return nil

			case <-ctx.Done():
				return ctx.Err()
			}
		}, func(error) {
			cancel()
			upg.Stop()
		}
}
