package main

import (
	"net"
	"sync"
)

type limitListener struct {
	net.Listener
	sem chan struct{}
	closeOnce sync.Once
	done chan struct{}
}

func (l *limitListener) release() {
	<-l.sem
}

func (l *limitListener) acquire() bool {
	select {
	case <-l.done:
		return false
	case l.sem <- struct{}{}:
		return true
	}
}

func (l *limitListener) Close() error {
	err := l.Listener.Close()
	l.closeOnce.Do(func() {
		close(l.done)
	})
	return err
}
 
func (l *limitListener) Accept() (net.Conn, error) {
	acquired := l.acquire()

	c, err := l.Listener.Accept()
	if err != nil {
		if acquired {
			l.release()
		}
		return nil, err
	}

	return &limitListenerConn{
		Conn: c,
		release:l.release,
	}, nil
}

type limitListenerConn struct {
	net.Conn
	releaseOnce sync.Once
	release func()
}

func (l *limitListenerConn) Close() error {
	err := l.Conn.Close()
	l.releaseOnce.Do(l.release)
	return err
}