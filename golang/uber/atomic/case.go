package main

import "go.uber.org/atomic"

type foo struct {
	running atomic.Bool
}

func (f *foo) start() {
	if f.running.Swap(true) {
		// already running…
		return
	}
	// start the Foo
}

func (f *foo) isRunning() bool {
	return f.running.Load()
}
