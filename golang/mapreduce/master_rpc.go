package mapreduce

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
)

func (m *Master) Shutdown(_, _ *struct{}) error {
	debug("Shutdown: registration server\n")
	close(m.shutdown)
	m.l.Close() // causes the Accept to fail
	return nil
}

func (m *Master) startRPCServer() {
	rpcs := rpc.NewServer()
	rpcs.Register(m)

	os.Remove(m.address)

	l, err := net.Listen("unix", m.address)
	if err != nil {
		log.Fatal("RegstrationServer", m.address, " error: ", err)
	}
	m.l = l

	go func() {
	loop:
		for {
			select {
			case <-m.shutdown:
				break loop
			default:
			}
			conn, err := m.l.Accept()
			if err == nil {
				go func() {
					rpcs.ServeConn(conn)
					conn.Close()
				}()
			} else {
				debug("RegistrationServer: accept error", err)
				break
			}
		}
		debug("RegistrationServer: done\n")
	}()
}

func (m *Master) stopRPCServer() {
	var reply ShutdownReply
	ok := call(m.address, "Master.Shutdown", new(struct{}), &reply)
	if ok == false {
		fmt.Printf("Cleanup: RPC %s error\n", m.address)
	}
	debug("cleanupRegistration: done\n")
}
