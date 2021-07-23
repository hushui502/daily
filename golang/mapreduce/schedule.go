package mapreduce

import "fmt"

func (m *Master) schedule(phase jobPhase) {
	var ntasks int
	// number of inputs (for reduce) or outputs (for map)
	var nios int

	switch phase {
	case mapPhase:
		ntasks = len(m.files)
		nios = m.nReduce
	case reducePhase:
		ntasks = m.nReduce
		nios = len(m.files)
	}

	fmt.Printf("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, nios)

	done := make(chan bool)
	for i := 0; i < ntasks; i++ {
		go func(number int) {
			args := DoTaskArgs{m.jobName, m.files[number], phase, number, nios}
			var worker string
			reply := new(struct{})
			ok := false
			for ok != true {
				worker = <-m.registerChannel
				ok = call(worker, "Worker.DoTask", args, reply)
			}
			done <- true
			m.registerChannel <- worker
		}(i)
	}

	for i := 0; i < ntasks; i++ {
		<-done
	}

	fmt.Printf("Schedule: %v phase done\n", phase)
}
