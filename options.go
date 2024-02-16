package run

import "go.uber.org/zap"

type Options struct {
	logger *zap.Logger
}

type Option func(*Options)

func WithLogger(logger *zap.Logger) Option {
	return func(o *Options) {
		o.logger = logger
	}
}
