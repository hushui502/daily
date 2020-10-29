package scheduler

import "awesomeProject2/engine"

type SimpleScheduler struct {
	workChan chan engine.Request
}

func (s *SimpleScheduler) Submit(request engine.Request) {
	go func() {
		s.workChan <- request
	}()
}

func (s *SimpleScheduler) ConfigureMasterWorkerChan(c chan engine.Request) {
	s.workChan = c
}
