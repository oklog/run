package run

import "log/slog"

type Option func(*Group)

func WithLogger(handler slog.Handler) Option {
	return func(o *Group) {
		o.logger = slog.New(handler)
	}
}
