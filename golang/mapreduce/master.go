package mapreduce

import (
	"fmt"
	"net"
	"sync"
)

// Master holds all the state that the master needs to keep track of.
type Master struct {
	sync.Mutex

	address         string
	registerChannel chan string
	doneChannel     chan bool
	workers         []string // producted by the mutex

	jobName string   // name of currently executing job
	files   []string // inout files
	nReduce int      // number of reduce partitions

	shutdown chan struct{}
	l        net.Listener
	stats    []int
}

// NewMaster initializes a new Map/Reduce Master
func NewMaster(master string) (mr *Master) {
	mr = new(Master)
	mr.address = master
	mr.shutdown = make(chan struct{})
	mr.registerChannel = make(chan string)
	mr.doneChannel = make(chan bool)
	return
}

// Register is an RPC method that is called by workers after they have started
// up to report that they are ready to receive tasks.
func (m *Master) Register(args *RegisterArgs, _ *struct{}) error {
	m.Lock()
	defer m.Unlock()

	debug("Register: worker %s\n", args.Worker)

	go func() {
		m.registerChannel <- args.Worker
	}()

	return nil
}

func Sequential(jobName string, files []string, nreduce int,
	mapF func(string, string) []KeyValue,
	reduceF func(string, []string) string) (m *Master) {
	m = NewMaster("master")
	go m.run(jobName, files, nreduce, func(phase jobPhase) {
		switch phase {
		case mapPhase:
			for i, f := range m.files {
				doMap(m.jobName, i, f, m.nReduce, mapF)
			}
		case reducePhase:
			for i := 0; i < m.nReduce; i++ {
				doReduce(m.jobName, i, len(m.files), reduceF)
			}
		}
	}, func() {
		m.stats = []int{len(files) + nreduce}
	})

	return
}

// Distributed schedules map and reduce tasks on workers that register with the
// master over RPC.
func Distributed(jobName string, files []string, nreduce int, master string) (m *Master) {
	m = NewMaster(master)
	m.startRPCServer()
	go m.run(jobName, files, nreduce, m.schedule, func() {
		m.stats = m.killWorkers()
		m.stopRPCServer()
	})
	return
}

// run executes a mapreduce job on the given number of mappers and reducers.
func (m *Master) run(jobName string, files []string, nreduce int,
	schedule func(phase jobPhase),
	finish func()) {
	m.jobName = jobName
	m.files = files
	m.nReduce = nreduce

	fmt.Printf("%s: Starting Map/Reduce task %s\n", m.address, m.jobName)

	schedule(mapPhase)
	schedule(reducePhase)
	finish()
	m.merge()

	fmt.Printf("%s: Map/Reduce task completed\n", m.address)

	m.doneChannel <- true
}

func (m *Master) Wait() {
	<-m.doneChannel
}

// killWorkers cleans up all workers by sending each one a Shutdown RPC.
// It also collects and returns the number of tasks each worker has performed.
func (m *Master) killWorkers() []int {
	m.Lock()
	defer m.Unlock()

	ntasks := make([]int, 0, len(m.workers))
	for _, w := range m.workers {
		debug("Master: shutdown worker %s\n", w)
		var reply ShutdownReply
		ok := call(w, "Worker.Shutdown", new(struct{}), &reply)
		if ok == false {
			fmt.Printf("Master: rpc %s shutdown error\n", w)
		} else {
			ntasks = append(ntasks, reply.Ntasks)
		}
	}

	return ntasks
}
