package simpleratelimit

import (
	"sync/atomic"
	"time"
)

type RateLimiter struct {
	rate uint64
	allowance uint64
	max uint64
	unit uint64
	lastCheck uint64
}

func New(rate int, per time.Duration) *RateLimiter {
	nano := uint64(per)
	if nano < 1 {
		nano = uint64(time.Second)
	}
	if rate < 1 {
		rate = 1
	}

	return &RateLimiter{
		rate:      uint64(rate),
		allowance: uint64(rate) * nano,
		max:       uint64(rate) * nano,
		unit:      nano,
		lastCheck: unixNano(),
	}
}

func (rl *RateLimiter) Limit() bool {
	now := unixNano()

	passed := now - atomic.SwapUint64(&rl.lastCheck, now)
	rate := atomic.LoadUint64(&rl.rate)
	current := atomic.AddUint64(&rl.allowance, passed*rate)

	if max := atomic.LoadUint64(&rl.max); current > max {
		atomic.AddUint64(&rl.allowance, max-current)
		current = max
	}

	if current < rl.unit {
		return true
	}

	atomic.AddUint64(&rl.allowance, -rl.unit)
	return false
}

func (rl *RateLimiter) UpdateRate(rate int) {
	atomic.StoreUint64(&rl.rate, uint64(rate))
	atomic.StoreUint64(&rl.max, uint64(rate)*rl.unit)
}

func (rl *RateLimiter) Undo() {
	current := atomic.AddUint64(&rl.allowance, rl.unit)

	if max := atomic.LoadUint64(&rl.max); current > max {
		atomic.AddUint64(&rl.allowance, max-current)
	}
}

func unixNano() uint64 {
	return uint64(time.Now().UnixNano())
}