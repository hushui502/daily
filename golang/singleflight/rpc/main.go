package main

import (
	"fmt"
	"net/http"
	"net/rpc"
	"time"
)

type (
	Arg struct {
		Caller int
	}
	Data struct {}
)


func (d *Data) GetData(arg *Arg, reply *string) error {
	fmt.Printf("request from client %d\n", arg.Caller)
	time.Sleep(1 * time.Second)
	*reply = "source data from rpcServer"
	return nil
}

func main() {
	d := new(Data)
	rpc.Register(d)
	rpc.HandleHTTP()
	fmt.Println("start rpc server")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}