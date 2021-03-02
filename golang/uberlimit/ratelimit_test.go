package uberlimit

import (
	"github.com/andres-erbsen/clock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
	"sync"
	"testing"
	"time"
)

type testRunner interface {
	createLimiter(int, ...Option) Limiter
	startTaking(rls ...Limiter)
	assertCountAt(d time.Duration, count int)
	afterFunc(d time.Duration, fn func())
}

type runnerImpl struct {
	t *testing.T

	clock       *clock.Mock
	constructor func(int, ...Option) Limiter
	count       atomic.Int32
	maxDuration time.Duration
	doneCh      chan struct{}
	wg          sync.WaitGroup
}

func runTest(t *testing.T, fn func(testRunner)) {
	impls := []struct {
		name        string
		constructor func(int, ...Option) Limiter
	}{
		{
			name: "mutex",
			constructor: func(rate int, opts ...Option) Limiter {
				return newMutexBased(rate, opts...)
			},
		},
		{
			name: "atomic",
			constructor: func(rate int, opts ...Option) Limiter {
				return newAtomicBased(rate, opts...)
			},
		},
	}

	for _, tt := range impls {
		t.Run(tt.name, func(t *testing.T) {
			r := runnerImpl{
				t:           t,
				clock:       clock.NewMock(),
				constructor: tt.constructor,
				doneCh:      make(chan struct{}),
			}
			defer close(r.doneCh)
			defer r.wg.Wait()

			fn(&r)
			r.clock.Add(r.maxDuration)
		})
	}
}

func (r *runnerImpl) createLimiter(rate int, opts ...Option) Limiter {
	opts = append(opts, WithClock(r.clock))

	return r.constructor(rate, opts...)
}

func (r *runnerImpl) startTaking(rls ...Limiter) {
	r.gowait(func() {
		for {
			for _, rl := range rls {
				rl.Take()
			}
			r.count.Inc()
			select {
			case <-r.doneCh:
				return
			default:
			}
		}
	})
}

func (r *runnerImpl) assertCountAt(d time.Duration, count int) {
	r.wg.Add(1)
	r.afterFunc(d, func() {
		assert.Equal(r.t, int32(count), r.count.Load(), "count not as expected")
		r.wg.Done()
	})
}

func (r *runnerImpl) afterFunc(d time.Duration, fn func()) {
	if d > r.maxDuration {
		r.maxDuration = d
	}

	r.gowait(func() {
		select {
		case <-r.doneCh:
			return
		case <-r.clock.After(d):
		}

		fn()
	})
}

func (r *runnerImpl) gowait(fn func()) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		wg.Done()
		fn()
	}()

	wg.Wait()
}

func TestUnlimited(t *testing.T) {
	now := time.Now()
	rl := NewUnlimited()
	for i := 0; i < 1000; i++ {
		rl.Take()
	}
	assert.Condition(t, func() bool { return time.Since(now) < 1*time.Millisecond }, "no artificial delay")
}

func TestRateLimiter(t *testing.T) {
	runTest(t, func(r testRunner) {
		rl := r.createLimiter(100, WithoutSlack)

		// Create copious counts concurrently.
		r.startTaking(rl)
		r.startTaking(rl)
		r.startTaking(rl)
		r.startTaking(rl)

		r.assertCountAt(1*time.Second, 100)
		r.assertCountAt(2*time.Second, 200)
		r.assertCountAt(3*time.Second, 300)
	})
}

func TestPer(t *testing.T) {
	runTest(t, func(r testRunner) {
		rl := r.createLimiter(7, WithoutSlack, Per(time.Minute))

		r.startTaking(rl)
		r.startTaking(rl)

		r.assertCountAt(1*time.Second, 1)
		r.assertCountAt(1*time.Minute, 8)
		r.assertCountAt(2*time.Minute, 15)
	})
}
