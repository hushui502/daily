package gobreaker

import (
	"errors"
	"fmt"
	"sync"
	"time"
)
// 三种状态：
// 熔断器关闭状态，可以正常手淫
// 熔断器开启状态，禁止异常，开始戒色
// 熔断器半开状态，可以看片，并验证是否可以手淫

// 四种状态转移：
// 在熔断器关闭状态，当连续的前列腺炎发作，并满足一定条件后，直接熔断器开启，禁止手淫
// 在熔断器开启状态，如果戒色超过了一定时间，将进入半开启状态，可以看片，并验证是否可以手淫
// 在熔断器半开状态，如果出现手淫引起前列腺炎发作，则再次进入开启状态
// 在熔断器半开状态，所有的手淫都没引起前列腺炎发作，则熔断器关闭。可以正常手淫

type State int

const (
	StateClosed State = iota
	StateHalfOpen
	StateOpen
)

var (
	// ErrTooManyRequests is returned when the CB state is half open and the requests count is over the cb maxRequests
	ErrTooManyRequests = errors.New("too many requests")
	// ErrOpenState is returned when the CB state is open
	ErrOpenState = errors.New("circuit breaker is open")
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateHalfOpen:
		return "half-open"
	case StateOpen:
		return "open"
	default:
		return fmt.Sprintf("unknown state: %d", s)
	}
}

type Counts struct {
	Requests uint32
	TotalSuccess uint32
	TotalFailures uint32
	ConsecutiveSuccess uint32
	ConsecutiveFailure uint32
}

func (c *Counts) onRequest() {
	c.Requests++
}

func (c *Counts) onSuccess() {
	c.TotalSuccess++
	c.ConsecutiveSuccess++
	c.ConsecutiveFailure = 0
}

func (c *Counts) onFailure() {
	c.TotalFailures++
	c.ConsecutiveFailure++
	c.ConsecutiveSuccess = 0
}

func (c *Counts) clear() {
	c.Requests = 0
	c.TotalSuccess = 0
	c.TotalFailures = 0
	c.ConsecutiveSuccess = 0
	c.ConsecutiveFailure = 0
}

type Settings struct {
	Name string
	MaxRequests uint32
	Interval time.Duration
	Timeout time.Duration
	ReadyToTrip func(counts Counts) bool
	OnStateChange func(name string, from State, to State)
}

type CircuitBreaker struct {
	name string
	maxRequest uint32
	interval time.Duration
	timeout time.Duration
	readyToTrip func(counts Counts) bool
	onStateChange func(name string, from State, to State)

	mutex sync.Mutex
	state State
	generation uint64
	counts Counts
	expiry time.Time
}

type TwoStepCircuitBreaker struct {
	cb *CircuitBreaker
}

const defaultInterval = time.Duration(0) * time.Second
const defaultTimeout = time.Duration(60) * time.Second

func defaultReadyToTrip(counts Counts) bool {
	return counts.ConsecutiveFailure > 5
}

func NewCircuitBreaker(st Settings) *CircuitBreaker {
	cb := new(CircuitBreaker)

	cb.name = st.Name
	cb.onStateChange = st.OnStateChange

	if st.MaxRequests == 0 {
		cb.maxRequest = -1
	} else {
		cb.maxRequest = st.MaxRequests
	}

	if st.Interval <= 0 {
		cb.interval = defaultInterval
	} else {
		cb.interval = st.Interval
	}

	if st.Timeout <= 0 {
		cb.timeout = defaultTimeout
	} else {
		cb.timeout = st.Timeout
	}

	if st.ReadyToTrip == nil {
		cb.readyToTrip = defaultReadyToTrip
	} else {
		cb.readyToTrip = st.ReadyToTrip
	}

	cb.toNewGeneration(time.Now())

	return cb
}

func NewTwoStepCircuitBreaker(st Settings) *TwoStepCircuitBreaker {
	return &TwoStepCircuitBreaker{
		cb:NewCircuitBreaker(st),
	}
}

func (cb *CircuitBreaker) Name() string {
	return cb.name
}

func (cb *CircuitBreaker) State() State {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	now := time.Now()
	state, _ := cb.currentState(now)

	return state
}

func (cb *CircuitBreaker) Execute(req func() (interface{}, error)) (interface{}, error) {
	generation, err := cb.beforeRequest()
	if err != nil {
		return nil, err
	}

	defer func() {
		e := recover()
		if e != nil {
			cb.afterRequest(generation, false)
			panic(e)
		}
	}()

	result, err := req()
	cb.afterRequest(generation, err == nil)

	return result, err
}

func (tscb *TwoStepCircuitBreaker) Name() string {
	return tscb.cb.name
}

func (tscb *TwoStepCircuitBreaker) State() State {
	return tscb.cb.state
}

func (tscb *TwoStepCircuitBreaker) Allow() (done func(success bool), err error) {
	generation, err := tscb.cb.beforeRequest()
	if err != nil {
		return nil, err
	}

	return func(success bool) {
		tscb.cb.afterRequest(generation, success)
	}, nil
}

func (cb *CircuitBreaker) beforeRequest() (uint64, error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	now := time.Now()
	state, generation := cb.currentState(now)

	if state == StateOpen {
		return generation, ErrOpenState
	} else if state == StateHalfOpen && cb.counts.Requests >= cb.maxRequest {
		return generation, ErrTooManyRequests
	}

	cb.counts.onRequest()

	return generation, nil
}

func (cb *CircuitBreaker) afterRequest(before uint64, success bool) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	now := time.Now()
	state, generation := cb.currentState(now)

	if generation != before {
		return
	}
	if success {
		cb.onSuccess(state, now)
	} else {
		cb.onFailure(state, now)
	}
}

func (cb *CircuitBreaker) onSuccess(state State, now time.Time) {
	switch state {
	case StateClosed:
		cb.counts.onSuccess()
	case StateHalfOpen:
		cb.counts.onSuccess()
		if cb.counts.ConsecutiveSuccess >= cb.maxRequest {
			cb.setState(StateClosed, now)
		}
	}
}

func (cb *CircuitBreaker) onFailure(state State, now time.Time) {
	switch state {
	case StateClosed:
		cb.counts.onFailure()
		if cb.readyToTrip(cb.counts) {
			cb.setState(StateClosed, now)
		}
	}
}

func (cb *CircuitBreaker) currentState(now time.Time) (State, uint64) {
	switch cb.state {
	case StateClosed:
		if !cb.expiry.IsZero() && cb.expiry.Before(now) {
			cb.toNewGeneration(now)
		}
	case StateOpen:
		if cb.expiry.Before(now) {
			cb.setState(StateHalfOpen, now)
		}
	}

	return cb.state, cb.generation
}

func (cb *CircuitBreaker) setState(state State, now time.Time) {
	if cb.state == state {
		return
	}

	prev := cb.state
	cb.state = state

	cb.toNewGeneration(now)

	if cb.onStateChange != nil {
		cb.onStateChange(cb.name, prev, state)
	}
}

func (cb *CircuitBreaker) toNewGeneration(now time.Time) {
	cb.generation++
	cb.counts.clear()

	var zero time.Time
	switch cb.state {
	case StateOpen:
		cb.expiry = now.Add(cb.timeout)
	case StateClosed:
		if cb.interval == 0 {
			cb.expiry = zero
		} else {
			cb.expiry = now.Add(cb.interval)
		}
	default:
		cb.expiry = zero
	}
}






















