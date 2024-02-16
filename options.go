package run

import "go.uber.org/zap"

type Option func(*Group)

func WithLogger(logger *zap.Logger) Option {
	return func(o *Group) {
		o.logger = logger
	}
}
