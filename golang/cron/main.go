package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"sync/atomic"
	"time"
)


type GreetingJob struct {
	Name string
}

func (g GreetingJob) Run() {
	fmt.Println("Hello ", g.Name)
}

type panicJob struct {
	count int
}

func (p *panicJob) Run() {
	p.count++
	if p.count == 1 {
		panic("-0-0-0-0-0-0-0-0-0-0")
	}

	fmt.Println("hello --- ")
}

type delayJob struct {
	count int
}

func (d *delayJob) Run() {
	time.Sleep(time.Second * 2)
	d.count++
	log.Printf("%d hello\n", d.count)
}

type skipJob struct {
	count int32
}

func (d *skipJob) Run() {
	atomic.AddInt32(&d.count, 1)
	log.Printf("%d hello\n", d.count)
	if atomic.LoadInt32(&d.count) == 1 {
		time.Sleep(time.Second*2)
	}
}

func main() {
	c := cron.New()
	c.AddJob("@every 1s", cron.NewChain(cron.Recover(cron.DefaultLogger), cron.SkipIfStillRunning(cron.DefaultLogger)).Then(&skipJob{}))

	c.Start()

	time.Sleep(time.Second*12)
}
