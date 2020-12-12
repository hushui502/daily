package adapter

import "testing"

func TestAliyunClient_CreateServer(t *testing.T) {
	var a ICreateServer = &AliyunClientAdapter{Client:AliyunClient{}}

	a.CreateServer(1.0, 2.0)
}

func TestAwsClientAdapter_CreateServer(t *testing.T) {
	// 确保 adapter 实现了目标接口
	var a ICreateServer = &AwsClientAdapter{
		Client: AWSClient{},
	}

	a.CreateServer(1.0, 2.0)
}
