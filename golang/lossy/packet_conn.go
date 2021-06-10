package lossy

import (
	"math/rand"
	"net"
	"sync"
	"time"
)

type packetConn struct {
	net.PacketConn
	minLatency        time.Duration
	maxLatency        time.Duration
	packetLossRate    float64
	writeDeadline     time.Time
	closed            bool
	mu                *sync.Mutex
	rand              *rand.Rand
	throttleMu        *sync.Mutex
	timeToWaitPerByte float64
	headerOverhead    int
}

func NewPacketConn(c net.PacketConn, bandwidth int, minLatency, maxLatency time.Duration, packetLossRate float64, headerOverhead int) net.PacketConn {
	var timeToWaitPerByte float64
	if bandwidth <= 0 {
		timeToWaitPerByte = 0
	} else {
		timeToWaitPerByte = float64(time.Second) / float64(bandwidth)
	}
	return &packetConn{
		PacketConn:        c,
		minLatency:        minLatency,
		maxLatency:        maxLatency,
		packetLossRate:    packetLossRate,
		writeDeadline:     time.Time{},
		closed:            false,
		mu:                &sync.Mutex{},
		rand:              rand.New(rand.NewSource(time.Now().UnixNano())),
		throttleMu:        &sync.Mutex{},
		timeToWaitPerByte: timeToWaitPerByte,
		headerOverhead:    headerOverhead,
	}
}

func (c *packetConn) WriteTo(p []byte, addr net.Addr) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed || !c.writeDeadline.Equal(time.Time{}) && c.writeDeadline.Before(time.Now()) {
		return c.PacketConn.WriteTo(p, addr)
	}
	go func() {
		c.throttleMu.Lock()
		time.Sleep(time.Duration(c.timeToWaitPerByte * (float64(len(p) + c.headerOverhead))))
		c.throttleMu.Unlock()
		if c.rand.Float64() >= c.packetLossRate {
			time.Sleep(c.minLatency + time.Duration(float64(c.maxLatency-c.minLatency)*c.rand.Float64()))
			c.mu.Lock()
			_, _ = c.PacketConn.WriteTo(p, addr)
			c.mu.Unlock()
		}
	}()
	return len(p), nil
}

func (c *packetConn) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed = true
	return c.PacketConn.Close()
}

func (c *packetConn) SetDeadline(t time.Time) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.writeDeadline = t
	return c.PacketConn.SetDeadline(t)
}

func (c *packetConn) SetWriteDeadline(t time.Time) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.writeDeadline = t
	return c.PacketConn.SetWriteDeadline(t)
}

