package ffcron

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	OnMode  = true
	OffMode = false
)

var (
	ErrNotFoundJob     = errors.New("not found job")
	ErrAlreadyRegister = errors.New("the job already in pool")
	ErrJobDoFuncNil    = errors.New("callback func is nil")
	ErrCronSpecInvalid = errors.New("crontab spec is invalid")
)

// null logger
var defaultLogger = func(level, s string) {}

type loggerType func(level, s string)

type panicType func(srv, err string)

var panicCaller = func(srv, err string) {}

func SetPanicCaller(p panicType) {
	panicCaller = p
}

type CronScheduler struct {
	tasks  map[string]*JobModel
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
	once   *sync.Once
	sync.RWMutex
}

func New() *CronScheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &CronScheduler{
		tasks:  make(map[string]*JobModel),
		ctx:    ctx,
		cancel: cancel,
		wg:     &sync.WaitGroup{},
		once:   &sync.Once{},
	}
}

func (c *CronScheduler) Register(srv, model *JobModel) error {

}

func (c *CronScheduler) reset(srv string, model *JobModel, denyReplace, autoStart bool) error {
	c.Lock()
	defer c.Unlock()

	// validate model by robfig/cron/parse
	err := model.validate()
	if err != nil {
		return err
	}

	model.srv = srv
	model.ctx = c.ctx

	oldModel, ok := c.tasks[srv]
	if denyReplace && ok {
		return ErrAlreadyRegister
	}

	if ok {
		// kill oldmodel for starting a new model to replace it
		oldModel.kill()
	}

	c.tasks[srv] = model
	if autoStart {
		c.wg.Add(1)
		go c.tasks[srv].runLoop(c.wg)
	}

	return nil
}

func (c *CronScheduler) UnRegister(srv string) error {
	c.Lock()
	defer c.Unlock()

	oldModel, ok := c.tasks[srv]
	if !ok {
		return ErrNotFoundJob
	}

	oldModel.kill()
	delete(c.tasks, srv)

	return nil
}

func (c *CronScheduler) Stop() {
	c.Lock()
	defer c.Unlock()

	for srv, job := range c.tasks {
		job.kill()
		delete(c.tasks, srv)
	}
}

func (c *CronScheduler) StopService(srv string) {
	c.Lock()
	defer c.Unlock()

	job, ok := c.tasks[srv]
	if !ok {
		return
	}
	job.kill()
	delete(c.tasks, srv)
}

func (c *CronScheduler) StopServicePrefix(regex string) {
	c.Lock()
	defer c.Unlock()

	for srv, job := range c.tasks {
		if !strings.Contains(srv, regex) {
			continue
		}
		job.kill()
		delete(c.tasks, srv)
	}
}

func validateSpec(spec string) bool {
	_, err := Parse(spec)
	if err != nil {
		return false
	}

	return true
}

func getNextDue(spec string) (time.Time, error) {
	sc, err := Parse(spec)
	if err != nil {
		return time.Now(), err
	}

	// simulated job
	time.Sleep(10 * time.Millisecond)
	due := sc.Next(time.Now())

	return due, err
}

func getNextDueSafe(spec string, last time.Time) (time.Time, error) {
	var (
		due time.Time
		err error
	)

	for {
		due, err = getNextDue(spec)
		if err != nil {
			return due, err
		}
		if last.Equal(due) {
			// avoid lost some accuracy
			time.Sleep(100 * time.Millisecond)
			continue
		}

		break
	}

	return due, nil
}

func (c *CronScheduler) Start() {
	c.once.Do(func() {
		for _, job := range c.tasks {
			c.wg.Add(1)
			job.runLoop(c.wg)
		}
	})
}

func (c *CronScheduler) Wait() {
	c.wg.Wait()
}

func (c *CronScheduler) GetServiceCron(srv string) (*JobModel, error) {
	c.RLock()
	defer c.RUnlock()

	oldModel, ok := c.tasks[srv]
	if !ok {
		return nil, ErrNotFoundJob
	}

	return oldModel, nil
}

func NewJobModel(spec string, f func(), options ...JobOption) (*JobModel, error) {
	var err error
	job := &JobModel{
		running:    true,
		async:      false,
		do:         f,
		spec:       spec,
		notifyChan: make(chan int, 1),
	}

	for _, opt := range options {
		if opt != nil {
			if err := opt(job); err != nil {
				return nil, err
			}
		}
	}

	err = job.validate()
	if err != nil {
		return nil, err
	}

	return job, nil
}

type JobModel struct {
	// srv name
	srv string
	// call func
	do func()
	// if async = true => go func() { do() }
	async bool
	// try catch panic
	tryCatch bool
	// cron spec
	spec       string
	ctx        context.Context
	notifyChan chan int
	// break for loop
	running bool
	// ensure job worker is exited already
	exited bool
	sync.RWMutex
}

type JobOption func(*JobModel) error

func AsynMode() JobOption {
	return func(o *JobModel) error {
		o.async = true
		return nil
	}
}

func TryCatchMode() JobOption {
	return func(o *JobModel) error {
		o.tryCatch = true
		return nil
	}
}

func (j *JobModel) SetTryCatch(b bool) {
	j.tryCatch = b
}

func (j *JobModel) SetAsynMode(b bool) {
	j.async = b
}

func (j *JobModel) validate() error {
	if j.do == nil {
		return ErrJobDoFuncNil
	}
	if _, err := getNextDue(j.spec); err != nil {
		return err
	}

	return nil
}

func (j *JobModel) runLoop(wg *sync.WaitGroup) {
	go j.run(wg)
}

func (j *JobModel) run(wg *sync.WaitGroup) {
	var (
		doTimeCostFunc = func() {
			startTS := time.Now()
			defaultLogger("info",
				fmt.Sprintf("scheduler service: %s begin run",
					j.srv))

			if j.tryCatch {
				tryCatch(j)
			} else {
				j.do()
			}

			defaultLogger("info",
				fmt.Sprintf("scheduler service: %s has been finished, time cost: %s spec: %s",
					j.srv,
					time.Since(startTS).String(),
					j.spec))
		}

		timer        *time.Timer
		lastNextTime time.Time
		due          time.Time
		interval     time.Duration

		err error
	)

	due, err = getNextDue(j.spec)
	interval = due.Sub(time.Now())
	if err != nil {
		panic(err.Error())
	}

	lastNextTime = due
	defaultLogger("info",
		fmt.Sprintf("scheduler service: %s next time is: %s, sub: %s",
			j.srv,
			due.String(),
			interval.String()))

	// init timer
	timer = time.NewTimer(interval)

	// clean
	defer func() {
		timer.Stop()
		wg.Done()
		j.exited = true
	}()

	for j.running {
		select {
		case <-timer.C:
			if time.Now().Before(due) {
				timer.Reset(
					due.Sub(time.Now()) + 50*time.Millisecond)
				continue
			}

			due, _ := getNextDueSafe(j.spec, lastNextTime)
			lastNextTime = due
			interval := due.Sub(time.Now())
			timer.Reset(interval)

			if j.async {
				go doTimeCostFunc()
			} else {
				doTimeCostFunc()
			}

			defaultLogger("info",
				fmt.Sprintf("scheduler service: %s next time is %s sub: %s",
					j.srv,
					due.String(),
					interval.String()))

		case <-j.notifyChan:
			continue
		case <-j.ctx.Done():
			return
		}
	}
}

func (j *JobModel) kill() {
	j.running = false
	close(j.notifyChan)
}

func (j *JobModel) workerExited() bool {
	return j.exited
}

func (j *JobModel) notifySig() {
	select {
	case j.notifyChan <- 1:
	default:
		return
	}
}

func tryCatch(job *JobModel) {
	defer func() {
		if err := recover(); err != nil {
			panicCaller(
				job.srv,
				fmt.Sprintf("%v", err))

			defaultLogger("error",
				fmt.Sprintf("srv: %s, trycatch panicing %v", job.srv, err))
		}
	}()

	job.do()
}
