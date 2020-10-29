// Create: 2018/07/16 15:13:00 Change: 2018/07/28 13:27:31
// FileName: date.go
// Copyright (C) 2018 lijiaocn <lijiaocn@foxmail.com>
//
// Distributed under terms of the GPL license.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"log"
)


//定义文件File结构体，用于调用"私密通讯"项目的链码
type File struct {
	CouchDbDoctype string `json:"docType"` // 文件表
	FileNo         string `json:"fileNo"`  // file number
	UserId         string `json:"userId"`  // ID of File
	//	CreateTraceId  []string `json:"createTraceId"` //创建轨迹id
	FileMD5 string `json:"fileMD5"` // file MD5
	//ChangeId       string   `json:"changeId"`      //变更id，转发路径id，阅读id...
	//ChangeType     string   `json:"changeType"`    //变更类型 转发，阅读...
	CreateTime string `json:"createTime"` //创建时间
}

func main() {

	//读取配置文件，创建SDK对象
	configProvider := config.FromFile("tls-off/fixtures/config-without-tls.yaml")
	sdk, err := fabsdk.New(configProvider)
	if err != nil {
		log.Fatalf("create sdk fail: %s\n", err.Error())
	}else{
		fmt.Println("create a new fabsdk")
	}


	//构建channel客户端,
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
	var file File
	file.FileNo="fileNo2"
	fileByte,err:=json.Marshal(file)

	var args [][]byte
	args = append(args, fileByte)

	request := channel.Request{
		ChaincodeID: "com",
		Fcn:         "queryFileCreate",
		Args:        args,
	}
	//调用channelClient执行请求
	response, err := channelClient.Execute(request)
	if err != nil {
		log.Fatal("query fail: ", err.Error())
	} else {
		fmt.Printf("response is %s\n", response.Payload)
	}
}
