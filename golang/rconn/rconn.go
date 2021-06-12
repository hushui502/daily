// github.com/jafagithubrlihi/rconn
package main

import (
	"net"
	"os"
)

var (
	remoteConn net.Conn
	localConn  net.Conn
)

func main() {
	initLogger()

	if len(os.Args) < 2 || (os.Args[1] != "-s" && os.Args[1] != "-c") {
		Log.Infof("Usage %s: [-s remote_port local_port | -c remote_addr remote_port local_addr local_port]", os.Args[0])
		os.Exit(1)
	}

	if os.Args[1] == "-s" {
		remoteListen, err := net.Listen("tcp", "0.0.0.0"+os.Args[2])
		if err != nil {
			Log.Errorf("Error listening on remote port: %s", err.Error())
			os.Exit(1)
		}
		defer remoteListen.Close()

		localListen, err := net.Listen("tcp", "0.0.0.0:"+os.Args[3])
		if err != nil {
			Log.Errorf("Error listening on local port: %s", err.Error())
			os.Exit(1)
		}
		defer localListen.Close()

		remoteConn, err = remoteListen.Accept()
		if err != nil {
			Log.Errorf("Error connecting on remote port: %s", err.Error())
			os.Exit(1)
		}

		localConn, err = localListen.Accept()
		if err != nil {
			Log.Errorf("Error connecting on local port: %s", err.Error())
			os.Exit(1)
		}

		statusRemote := make(chan bool)
		statusLocal := make(chan bool)

		go pipeSocket(true, statusRemote)
		go pipeSocket(false, statusLocal)

		for {
			status := <-statusLocal
			if !status {
				localConn, err = localListen.Accept()
				if err != nil {
					Log.Errorf("Error connecting on local port: %s", err.Error())
					os.Exit(1)
				}
			}
			go pipeSocket(false, statusLocal)
		}
	}

	if os.Args[1] == "-c" {
		var err error

		remoteConn, err = net.Dial("tcp", os.Args[2]+":"+os.Args[3])
		if err != nil {
			Log.Errorf("Error dialing to remote port: %s", err.Error())
			os.Exit(1)
		}

		localConn, err = net.Dial("tcp", os.Args[4]+":"+os.Args[5])
		if err != nil {
			Log.Errorf("Error dialing to local port: %s", err.Error())
			os.Exit(1)
		}

		statusRemote := make(chan bool)
		statusLocal := make(chan bool)

		go pipeSocket(true, statusRemote)
		go pipeSocket(false, statusLocal)

		for {
			status := <-statusLocal
			if !status {
				localConn, err = net.Dial("tcp", os.Args[4]+":"+os.Args[5])
				if err != nil {
					Log.Errorf("Error dialing to local port: %s", err.Error())
					os.Exit(1)
				}
			}
			go pipeSocket(false, statusLocal)
		}
	}
}

func pipeSocket(remoteToLocal bool, status chan<- bool) {
	for {
		buf := make([]byte, 1024)
		var err error
		var read int

		if remoteToLocal {
			read, err = remoteConn.Read(buf)
		} else {
			read, err = localConn.Read(buf)
		}
		if err != nil {
			Log.Errorf("Read error: %s", err.Error())
			status <- false
			return
		}

		if remoteToLocal {
			_, err = localConn.Write(buf[:read])
		} else {
			_, err = remoteConn.Write(buf[:read])
		}

		if err != nil {
			Log.Errorf("Write error: %s", err.Error())
			os.Exit(1)
		}
	}
}
