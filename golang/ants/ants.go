package ants

import (
	"errors"
	"log"
	"math"
	"os"
	"runtime"
	"time"
)

const (
	// DefaultAntsPoolSize is the default capacity for a default goroutine pool.
	DefaultAntsPoolSize = math.MaxInt32

	// DefaultCleanIntervalTime is the interval time to clean up goroutines.
	DefaultCleanIntervalTime = time.Second
)

const (
	// OPENED represents that the pool is opened.
	OPENED = iota

	// CLOSED represents that the pool is closed.
	CLOSED
)

var (
	// Error types for the Ants API.
	//---------------------------------------------------------------------------

	// ErrInvalidPoolSize will be returned when setting a negative number as pool capacity, this error will be only used
	// by pool with func because pool without func can be infinite by setting up a negative capacity.
	ErrInvalidPoolSize = errors.New("invalid size for pool")

	// ErrLackPoolFunc will be returned when invokers don't provide function for pool.
	ErrLackPoolFunc = errors.New("must provide function for pool")

	// ErrInvalidPoolExpiry will be returned when setting a negative number as the periodic duration to purge goroutines.
	ErrInvalidPoolExpiry = errors.New("invalid expiry for pool")

	// ErrPoolClosed will be returned when submitting task to a closed pool.
	ErrPoolClosed = errors.New("this pool has been closed")

	// ErrPoolOverload will be returned when the pool is full and no workers available.
	ErrPoolOverload = errors.New("too many goroutines blocked on submit or Nonblocking is set")

	// ErrInvalidPreAllocSize will be returned when trying to set up a negative capacity under PreAlloc mode.
	ErrInvalidPreAllocSize = errors.New("can not set up a negative capacity under PreAlloc mode")

	//---------------------------------------------------------------------------

	// workerChanCap determines whether the channel of a worker should be a buffered channel
	// to get the best performance. Inspired by fasthttp at
	// https://github.com/valyala/fasthttp/blob/master/workerpool.go#L139
	workerChanCap = func() int {
		if runtime.GOMAXPROCS(0) == 1 {
			return 0
		}

		return 1
	}()

	defaultLogger = Logger(log.New(os.Stderr, "", log.LstdFlags))

	defaultAntsPool, _ = NewPool(DefaultAntsPoolSize)

)

type Logger interface {
	Printf(format string, args ...interface{})
}

func Submit(task func()) error {
	return defaultAntsPool.Submit(task)
}

func Running() int {
	return defaultAntsPool.Running()
}


// Cap returns the capacity of this default pool.
func Cap() int {
	return defaultAntsPool.Cap()
}

// Free returns the available goroutines to work.
func Free() int {
	return defaultAntsPool.Free()
}

// Release Closes the default pool.
func Release() {
	defaultAntsPool.Release()
}

// Reboot reboots the default pool.
func Reboot() {
	defaultAntsPool.Reboot()
}

