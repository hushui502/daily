package connection

import (
	"bytes"
	"godis/lib/sync/wait"
	"net"
	"sync"
	"time"
)

type Connection struct {
	conn net.Conn
	// waiting until reply finished
	waitingReply wait.Wait
	// lock while server sending response
	mu   sync.Mutex
	subs map[string]bool
	// password may be changed by CONFIG command during runtime, so store the password
	password string
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Connection) Close() error {
	c.waitingReply.WaitWithTimeout(10 * time.Second)
	err := c.conn.Close()

	return err
}

func NewConn(conn net.Conn) *Connection {
	return &Connection{
		conn: conn,
	}
}

func (c *Connection) Write(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	_, err := c.conn.Write(b)

	return err
}

// Subscribe add current connections into subscribers of the given channel
func (c *Connection) Subscribe(channel string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.subs == nil {
		c.subs = make(map[string]bool)
	}
	c.subs[channel] = true
}

func (c *Connection) UnSubscribe(channel string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.subs == nil {
		return
	}
	delete(c.subs, channel)
}

func (c *Connection) SubsCount() int {
	if c.subs == nil {
		return 0
	}
	return len(c.subs)
}

func (c *Connection) GetChannels() []string {
	if len(c.subs) == 0 {
		return []string{}
	}
	channels := make([]string, len(c.subs))

	i := 0
	for sub := range c.subs {
		channels[i] = sub
		i++
	}

	return channels
}

// SetPassword stores password for authentication
func (c *Connection) SetPassword(password string) {
	c.password = password
}

// GetPassword get password for authentication
func (c *Connection) GetPassword() string {
	return c.password
}

// FakeConn implements redis.Connection for test
type FakeConn struct {
	Connection
	buf bytes.Buffer
}

// Write writes data to buffer
func (c *FakeConn) Write(b []byte) error {
	c.buf.Write(b)
	return nil
}

// Clean resets the buffer
func (c *FakeConn) Clean() {
	c.buf.Reset()
}

// Bytes returns written data
func (c *FakeConn) Bytes() []byte {
	return c.buf.Bytes()
}
