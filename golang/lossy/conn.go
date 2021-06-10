package lossy

import (
	"math/rand"
	"net"
	"sync"
	"time"
)

type conn struct {
	net.Conn
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

func NewConn(c net.Conn, bandwidth int, minLatency, maxLatency time.Duration, packetLossRate float64, headerOverhead int) net.Conn {
	var timeToWaitPerByte float64
	if bandwidth <= 0 {
		timeToWaitPerByte = 0
	} else {
		timeToWaitPerByte = float64(time.Second) / float64(bandwidth)
	}

	return &conn{
		Conn:              c,
		minLatency:        minLatency,
		maxLatency:        maxLatency,
		packetLossRate:    packetLossRate,
		writeDeadline: time.Time{},
		closed: false,
		mu: &sync.Mutex{},
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
		throttleMu: &sync.Mutex{},
		timeToWaitPerByte: timeToWaitPerByte,
		headerOverhead: headerOverhead,
	}
}

func (c *conn) Write(b []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed || !c.writeDeadline.Equal(time.Time{}) && c.writeDeadline.Before(time.Now()) {
		return c.Conn.Write(b)
	}

	go func() {
		c.throttleMu.Lock()
		time.Sleep(time.Duration(c.timeToWaitPerByte * (float64(len(b) + c.headerOverhead))))
		c.throttleMu.Unlock()
		if c.rand.Float64() >= c.packetLossRate {
			time.Sleep(c.minLatency + time.Duration(float64(c.maxLatency-c.minLatency)*c.rand.Float64()))
			c.mu.Lock()
			_, _ = c.Conn.Write(b)
			c.mu.Unlock()
		}
	}()

	return len(b), nil
}

func (c *conn) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed = true
	return c.Conn.Close()
}

func (c *conn) SetDeadline(t time.Time) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.writeDeadline = t
	return c.Conn.SetDeadline(t)
}

func (c *conn) SetWriteDeadline(t time.Time) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.writeDeadline = t
	return c.Conn.SetWriteDeadline(t)
}