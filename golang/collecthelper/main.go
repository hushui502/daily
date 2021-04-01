package collecthelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"sync"
	"time"

)

// These values may be changed to configure the thresholds observed by Check.
var (
	OversleepThreshold = 10 * time.Microsecond
	ChanSendThreshold  = 10 * time.Microsecond
	PingPongThreshold  = 10 * time.Microsecond
	ChainThreshold     = 100 * time.Microsecond
)

// Wraner is anything that can log warnings.
type Warner interface {
	Warningf(string, ...interface{})
}

func Check(w Warner) {
	checkChan <-w
}

const (
	sampleInterval   = 1 * time.Second
	testSleep        = 50 * time.Millisecond
	historySize      = 100
	numChainRoutines = 20
)

var (
	mu        sync.Mutex
	nextIndex int
	samples   [historySize]sample
)

type sample struct {
	start     time.Time
	oversleep time.Duration // undesired extra sleep latency
	bufSend   time.Duration // send on a buffered channel
	pingPong  time.Duration // ping-pong with goroutine on buffered channel
	chain     time.Duration
}

func init() {
	head = make(chan bool)
	tail = head
	for i := 0; i < numChainRoutines; i++ {
		ch := make(chan bool)
		go func(a, b chan bool) {
			for {
				b <- <-a
			}
		}(tail, ch)
		tail = ch
	}

	go channelHelper()
	go collectSampleLoop()
}

var (
	unbufc     = make(chan bool)
	bufc       = make(chan bool, 1)
	head, tail chan bool
)

func collectSampleLoop() {
	ticker := time.NewTicker(sampleInterval - testSleep)
	bad := true
	for {
		select {
		case <-ticker.C:
			s := collectSample()
			if overThreshold(&s) {
				bad = true
			}
		case w := <-checkChan:
			if bad {
				w.Warningf("Recent sample exceeded threshold.\nLast %v samples:\n%s", historySize, Samples())
				w.Warningf("Memory statistics:\n%s", memStats())
				bad = false
			}
			
		}
	}
}

func memStats() []byte {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	b, _ := json.MarshalIndent(stats, "", "	")
	return b
}

func overThreshold(s *sample) bool {
	return s.oversleep > OversleepThreshold ||
			s.bufSend > ChanSendThreshold ||
			s.pingPong > PingPongThreshold ||
			s.chain > ChainThreshold
}

var checkChan = make(chan Warner)

func channelHelper() {
	for {
		unbufc <- <-bufc
	}
}

func collectSample() sample {
	var s sample

	s.start = time.Now()
	time.Sleep(testSleep)
	t1 := time.Now()
	s.oversleep = t1.Sub(s.start) - testSleep

	bufc <- true
	t2 := time.Now()
	s.bufSend = t2.Sub(t1)
	<-unbufc
	t3 := time.Now()
	s.pingPong = t3.Sub(t2)

	head <- true
	<-tail
	s.chain = time.Now().Sub(t3)

	mu.Lock()
	defer mu.Unlock()
	idx := nextIndex
	nextIndex = (nextIndex + 1) % historySize
	samples[idx] = s

	return s
}

const header = "| " +
	"Sampled at | " +
	"Oversleep  | " +
	"Chan send  | " +
	"Ping-pong  | " +
	"Chain      |"

func Samples() string {
	mu.Lock()
	defer mu.Unlock()

	var buf bytes.Buffer

	fmt.Fprintf(&buf, header)
	idx := nextIndex
	now := time.Now()

	for n := 0; n < historySize; n++ {
		idx--
		if idx < 0 {
			idx = historySize - 1
		}
		s := &samples[idx]
		if s.start.IsZero() {
			break
		}
		hl := ""
		if overThreshold(s) {
			hl = "<---"
		}
		fmt.Fprintf(&buf, "| %5.1fs ago | %10v | %10v | %10v | %10v |%s\n",
			now.Sub(s.start).Seconds(),
			s.oversleep, s.bufSend, s.pingPong, s.chain,
			hl)
	}

	return buf.String()
}

