package rpcdemo

import (
	"log"
	"net"
	"net/rpc"
)

type HelloService struct {
}

func (p *HelloService) Hello(request string, reply *string) error {
	*reply = "hello" + request
	return nil
}

func main() {
	RegisterHelloService(new(HelloService))
	listener, err := net.Listen("tcp", "localhost:9999")
	if err != nil {
		log.Fatal("listenTCP error", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {

		}
		go rpc.ServeConn(conn)
	}

}
