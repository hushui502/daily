// Create: 2018/07/16 15:13:00 Change: 2018/07/28 13:27:31
// FileName: date.go
// Copyright (C) 2018 lijiaocn <lijiaocn@foxmail.com>
//
// Distributed under terms of the GPL license.

package main

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"log"
)

func main() {

	//读取配置文件，创建SDK对象
	configProvider := config.FromFile("D:\\project\\go\\src\\awesomeProject2\\fabric\\go-sdk\\example-old\\tls-on\\fixtures\\config.yaml")
	sdk, err := fabsdk.New(configProvider)
	if err != nil {
		log.Fatalf("create sdk fail: %s\n", err.Error())
	}else{
		fmt.Println("create a new fabsdk")
	}

	//调用合约
	channelProvider := sdk.ChannelContext("mychannel",
		fabsdk.WithUser("Admin"),
		fabsdk.WithOrg("org1"))

	channelClient, err := channel.New(channelProvider)
	if err != nil {
		log.Fatalf("create channel client fail: %s\n", err.Error())
	} else {
		fmt.Println("create channelClient succeed")
	}

	//构建函数请求
	var args1 [][]byte
	args1 = append(args1, []byte("b"))
	args1 = append(args1, []byte("a"))
	args1 = append(args1, []byte("3"))
	request := channel.Request{
		ChaincodeID: "mycc",
		Fcn:         "invoke",
		Args:        args1,
	}
	//调用channelClient执行请求
	response, err := channelClient.Execute(request)
	if err != nil {
		log.Fatal("query fail: ", err.Error())
	} else {
		fmt.Printf("response is %s\n", response.Payload)
	}

	//构建函数请求
	var args [][]byte
	args = append(args, []byte("a"))

	request = channel.Request{
		ChaincodeID: "mycc",
		Fcn:         "query",
		Args:        args,
	}
	//调用channelClient执行请求
	response, err = channelClient.Query(request)
	if err != nil {
		log.Fatal("query fail: ", err.Error())
	} else {
		fmt.Printf("response is %s\n", response.Payload)
	}
}
