package main

import "errors"

type cgroupLimiter struct {
}

func (r *cgroupLimiter) free() error {
	return nil
}

func (r *cgroupLimiter) configure(pid int, core float64, mb int) error {
	return errors.New("don't support cgroup")
}
