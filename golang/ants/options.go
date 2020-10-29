package ants

import "time"

type Option func(opts *Options)

func loadOptions(options ...Option) *Options {
	opts := new(Options)
	for _, option := range options {
		option(opts)
	}

	return opts
}

type Options struct {
	ExpiryDuration time.Duration

	PreAlloc bool

	MaxBlockingTasks int

	NonBlocking bool

	PanicHandler func(interface{})

	Logger Logger
}

// WithOptions accepts the whole options config.
func WithOptions(options Options) Option {
	return func(opts *Options) {
		*opts = options
	}
}

// WithExpiryDuration sets up the interval time of cleaning up goroutines.
func WithExpiryDuration(expiryDuration time.Duration) Option {
	return func(opts *Options) {
		opts.ExpiryDuration = expiryDuration
	}
}

// WithPreAlloc indicates whether it should malloc for workers.
func WithPreAlloc(preAlloc bool) Option {
	return func(opts *Options) {
		opts.PreAlloc = preAlloc
	}
}

// WithMaxBlockingTasks sets up the maximum number of goroutines that are blocked when it reaches the capacity of pool.
func WithMaxBlockingTasks(maxBlockingTasks int) Option {
	return func(opts *Options) {
		opts.MaxBlockingTasks = maxBlockingTasks
	}
}

// WithNonblocking indicates that pool will return nil when there is no available workers.
func WithNonblocking(nonBlocking bool) Option {
	return func(opts *Options) {
		opts.NonBlocking = nonBlocking
	}
}

// WithPanicHandler sets up panic handler.
func WithPanicHandler(panicHandler func(interface{})) Option {
	return func(opts *Options) {
		opts.PanicHandler = panicHandler
	}
}

// WithLogger sets up a customized logger.
func WithLogger(logger Logger) Option {
	return func(opts *Options) {
		opts.Logger = logger
	}
}
