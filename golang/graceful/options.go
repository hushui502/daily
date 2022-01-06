package graceful

import "context"

type Option func(*options)

type options struct {
	ctx context.Context
	logger Logger
}

func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

func WithLogger(logger Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

func NewOption(opts ...Option) options {
	defaultOptions := options{
		ctx: context.Background(),
		logger: NewLogger(),
	}

	for _, opt := range opts {
		opt(&defaultOptions)
	}

	return defaultOptions
}