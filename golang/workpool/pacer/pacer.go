package pacer

import "time"

type Pacer struct {
	delay  time.Duration
	gate   chan struct{}
	pause  chan struct{}
	paused chan struct{}
}

func NewPacer(delay time.Duration) *Pacer {
	p := &Pacer{
		delay:  delay,
		gate:   make(chan struct{}),
		pause:  make(chan struct{}, 1),
		paused: make(chan struct{}, 1),
	}

	go p.run()
	return p
}

func (p *Pacer) Pace(task func()) func() {
	return func() {
		p.Next()
		task()
	}
}

func (p *Pacer) Next() {
	// wait for item to be read from gate.
	p.gate <- struct{}{}
}

func (p *Pacer) Stop() {
	close(p.gate)
}

func (p *Pacer) IsPaused() bool {
	return len(p.paused) != 0
}

func (p *Pacer) Pause() {
	p.pause <- struct{}{}
	p.paused <- struct{}{}
}

func (p *Pacer) Resume() {
	<-p.pause
	<-p.paused
}

func (p *Pacer) run() {
	for range p.gate {
		time.Sleep(p.delay)
		// will wait here if channel blocked
		p.pause <- struct{}{}
		// clear channel
		<-p.pause
	}
}
