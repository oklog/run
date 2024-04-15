package run

import (
	"log/slog"
	"time"
)

type Option func(*Group)

// WithLogger is a functional option for setting the logger.
func WithLogger(logger *slog.Logger) Option {
	return func(o *Group) {
		o.logger = logger
	}
}

// WithCloseTimeout
func WithCloseTimeout(duration time.Duration) Option {
	return func(o *Group) {
		o.closeTimeout = duration
	}
}

// WithSyncShutdown ensures that a Runnable's Close method
// returns before shutting down the next Runnable.
func WithSyncShutdown() Option {
	return func(o *Group) {
		o.syncShutdown = true
	}
}
