package main

import (
	"fmt"
)

const (
	ALPHA = 0.8
	BETA  = 1.6
)

const (
	LBOUND = 100
	UBOUND = 60000
)

type TCP struct {
	srtt int32
	rto  uint32
}

func main() {
	tcp := &TCP{}
	lastRTO := uint32(0)
	rtts := []int32{
		100, 103, 105, 102, 100, 103, 105, 100, 103, 105, 102, 100, 102, 105, 100, 103, 105, 104, 100, 103, 105,
		200, 203, 205, 202, 200, 203, 205, 200, 203, 205, 202, 200, 202, 205, 200, 203, 205, 204, 200, 203, 205,
		400, 403, 405, 404,
		100, 103, 105, 102, 100, 103, 105, 100, 103, 105, 102, 100, 102, 105, 100, 103, 105, 104, 100, 103, 105,
		// 400, 403, 405, 400, 403, 405, 404, 400, 404, 405, 400, 403, 405, 404, 400, 403, 405,
	}
	for i, rtt := range rtts {
		if lastRTO == 0 {
			lastRTO = uint32(rtt)
		}
		rto := tcp.updateRTT(rtt)
		if uint32(rtt) > lastRTO {
			fmt.Printf("warn: loss packet, index: %d\n", i-1)
		}
		fmt.Printf("index: %d, rtt: %d, rto: %d\n", i, rtt, rto)
		lastRTO = rto
	}
}

func (t *TCP) updateRTT(rtt int32) uint32 {
	var rto uint32
	if t.srtt == 0 {
		t.srtt = rtt
		rto = uint32(rtt)
	} else {
		t.srtt = int32(ALPHA*float32(t.srtt) + ((1 - ALPHA) * float32(rtt)))
		rto = min(UBOUND, max(LBOUND, uint32(BETA*float32(t.srtt))))
	}
	return rto
}

func min(a, b uint32) uint32 {
	if a <= b {
		return a
	}
	return b
}

func max(a, b uint32) uint32 {
	if a >= b {
		return a
	}
	return b
}
