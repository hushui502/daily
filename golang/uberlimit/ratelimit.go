package uberlimit

import (
	"github.com/andres-erbsen/clock"
	"time"
)

type Limiter interface {
	Take() time.Time
}

type Clock interface {
	Now() time.Time
	Sleep(time.Duration)
}

type config struct {
	clock Clock
	slack int
	per time.Duration
}

func New(rate int, opts ...Option) Limiter {
	return newAtomicBased(rate, opts...)
}

func buildConfig(opts []Option) config {
	c := config{
		clock: clock.New(),
		slack: 10,
		per:   time.Second,
	}

	for _, opt := range opts {
		opt.apply(&c)
	}

	return c
}

type Option interface {
	apply(*config)
}

type clockOption struct {
	clock Clock
}

func (o clockOption) apply(c *config) {
	c.clock = o.clock
}

func WithClock(clock Clock) Option {
	return clockOption{clock: clock}
}

type slackOption int

func (o slackOption) apply(c *config) {
	c.slack = int(o)
}

var WithoutSlack Option = slackOption(0)

type perOption time.Duration

func (p perOption) apply(c *config) {
	c.per = time.Duration(p)
}

func Per(per time.Duration) Option {
	return perOption(per)
}

type unlimited struct {}

func NewUnlimited() Limiter {
	return unlimited{}
}

func (unlimited) Take() time.Time {
	return time.Now()
}

