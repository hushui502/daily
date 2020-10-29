package server

import (
	"flag"
	"fmt"
	"io"
	"net"
	"runtime"
	"strings"
	"time"
)

var (
	localPort int
	remotePort int
)

func init() {
	flag.IntVar(&localPort, "l", 5200, "the user link port")
	flag.IntVar(&remotePort, "r", 3333, "client listen port")
}

type client struct {
	conn net.Conn
	read chan []byte
	write chan []byte
	exit chan error
	reConn chan bool
}

func (c *client) Read() {
	_ = c.conn.SetReadDeadline(time.Now().Add(time.Second*10))

	for {
		data := make([]byte, 10240)
		n, err := c.conn.Read(data)
		if err != nil && err != io.EOF {
			if strings.Contains(err.Error(), "timeout") {
				_ = c.conn.SetReadDeadline(time.Now().Add(time.Second*3))
				_, err = c.conn.Write([]byte("pi"))
				continue
			}
			fmt.Println("读取出现错误。。。。")
			c.exit <- err
		}

		if data[0] == 'p' && data[1] == 'i' {
			fmt.Println("server 收到心跳包")
			continue
		}
		c.read <- data[:n]
	}
}

func (c *client) Write() {
	for {
		select {
		case data := <-c.write:
			_, err := c.conn.Write(data)
			if err != nil && err != io.EOF {
				c.exit <- err
			}
		}
	}
}

type user struct {
	conn net.Conn
	read chan []byte
	write chan []byte
	exit chan error
}

func (u *user) Read() {
	_ = u.conn.SetReadDeadline(time.Now().Add(time.Second*200))
	for {
		data := make([]byte, 10240)
		n, err := u.conn.Read(data)
		if err != nil && err != io.EOF {
			u.exit <- err
		}

		u.read <- data[:n]
	}
}

func (u *user) Write() {
	for {
		select {
		case data := <-u.write:
			_, err := u.conn.Write(data)
			if err != nil && err != io.EOF {
				u.exit <- err
			}
		}
	}
}




func HandleClient(client *client, userConnChan chan net.Conn) {
	go client.Read()
	go client.Write()

	for {
		select {
		case err := <-client.exit:
			fmt.Printf("client 出现错误,开始重试, err : %s\n", err.Error())
			client.reConn <- true
			runtime.Goexit()
		case useConn := <-userConnChan:
			user := user{
			conn:useConn,
			read:make(chan []byte),
			write:make(chan []byte),
			exit:make(chan error),
			}

			go user.Read()
			go user.Write()

			go Handle(client, user)
		}
	}
}

func Handle(client *client, user *user) {
	for {
		select {
		case userRecv := <-user.read:
			client.write <- userRecv
		case clientRecv := <-client.read:
			user.write <- clientRecv
		case err := <-client.exit:
			fmt.Println("client出现错误, 关闭连接", err.Error())
			_ = client.conn.Close()
			_ = user.conn.Close()
			client.reConn <- true
			runtime.Goexit()
		case err := <-user.exit:
			fmt.Println("user出现错误，关闭连接", err.Error())
			_ = user.conn.Close()
		}
	}
}

func AcceptUserConn(userListener net.Listener, connChan chan net.Conn) {
	userConn, err := userListener.Accept()
	if err != nil {

	}
	connChan <-userConn
}

func main() {
	flag.Parse()

	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()

	clientListener, err := net.Listen("tcp", fmt.Sprintf(":%d", remotePort))
	if err != nil {

	}
	userListener, err := net.Listen("tcp", fmt.Sprintf(":%d", localPort))

	for {
		clientConn, _ := clientListener.Accept()
		client := &client{
			conn:   clientConn,
			read:   make(chan []byte),
			write:  make(chan []byte),
			exit:   make(chan error),
			reConn: make(chan bool),
		}

		userConnChan := make(chan net.Conn)
		go AcceptUserConn(userListener, userConnChan)
		go HandleClient(client, userConnChan)

		<-client.reConn
	}
}