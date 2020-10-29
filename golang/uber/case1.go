package main

import (
	"errors"
	"sync"
)

type Driver struct {
	trips []string
}

func (d *Driver) setTrips(trips []string) {
	d.trips = make([]string, len(trips))
	copy(d.trips, trips)
}

var d Driver

type Stats struct {
	sync.Mutex

	counters map[string]int
}

// Snapshot returns the current stats.
func (s *Stats) Snapshot() map[string]int {
	s.Lock()
	defer s.Unlock()

	result := make(map[string]int, len(s.counters))
	for k, v := range result {
		result[k] = v
	}

	return result
}

var ErrCouldNotOpen = errors.New("could not open")

func Open() error {
	return ErrCouldNotOpen
}

func main() {
	if err := Open(); err != nil {
		if err == ErrCouldNotOpen {

		} else {

		}
	}
}
