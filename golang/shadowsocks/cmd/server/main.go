package main

import (
	"fmt"
	"github.com/phayes/freeport"
	"log"
	"net"
	"os"
	"shadowsocks/cmd"
	"shadowsocks/core"
	"shadowsocks/server"
	"strconv"
)

var version = "master"

func main() {
	log.SetFlags(log.Lshortfile)

	// 优先从环境变量获取sever port
	port, err := strconv.Atoi(os.Getenv("SHADOWSOCKS_SERVER_PORT"))
	// 随机获取一个空闲port
	if err != nil {
		port, err = freeport.GetFreePort()
	}
	if err != nil {
		// 强制一个端口，可能会失败
		port = 7478
	}

	// 默认配置
	config := &cmd.Config{
		ListenAddr: fmt.Sprintf(":%d", port),
		// 密码随机生成
		Password: core.RandPassword(),
	}
	config.ReadConfig()
	config.SaveConfig()

	// 启动 server 端并监听
	lsServer, err := server.NewLsServer(config.Password, config.ListenAddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(lsServer.Listen(func(listenAddr *net.TCPAddr) {
		log.Println(fmt.Sprintf(`
lightsocks-server:%s 启动成功，配置如下：
服务监听地址：
%s
密码：
%s`, version, listenAddr, config.Password))
	}))
}
