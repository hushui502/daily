package leakybucket

import (
	"errors"
	"time"
)

var (
	// ErrorFull is returned when the amount requested to add exceeds the remaining space in the bucket.
	ErrorFull = errors.New("add exceeds free capacity")
)

type BucketI interface {
	// Capacity of the bucket.
	Capacity() uint

	// Remaining space in the bucket.
	Remaining() uint

	// Reset returns when the bucket will be drained.
	Reset() time.Time

	// Add to the bucket. Return bucket state after adding.
	Add(uint) (BucketState, error)
}

type BucketState struct {
	Capacity uint
	Remaining uint
	Reset time.Time
}

type StorageI interface {
	// Create a bucket with a name, capacity, and rate.
	// rate is how long it takes for full capacity to drain
	Create(name string, capacity uint, rate time.Duration) (BucketI, error)
}
