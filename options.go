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

// WithOrderedShutdown
func WithOrderedShutdown() Option {
	return func(o *Group) {
		o.orderedShutdown = true
	}
}
