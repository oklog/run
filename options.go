package run

import "log/slog"

type Option func(*Group)

// WithLogger is a functional option for setting the logger.
func WithLogger(logger *slog.Logger) Option {
	return func(o *Group) {
		o.logger = logger
	}
}
