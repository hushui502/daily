package local

import (
	"log"
	"net"
	"shadowsocks/core"
)

/*
	运行在本机的 local 端的职责是把本机程序发送给它的数据经过加密后转发给墙外的代理服务器，总体工作流程如下：

	监听来自本机浏览器的代理请求；
	转发前加密数据；
	转发socket数据到墙外代理服务端；
	把服务端返回的数据转发给用户的浏览器。
*/

type LsLocal struct {
	*core.SecureSocket
}

func New(password *core.Password, listenAddr, remoteAddr *net.TCPAddr) *LsLocal  {
	return &LsLocal{
		SecureSocket: &core.SecureSocket{
			Cipher:     core.NewCipher(password),
			ListenAddr: listenAddr,
			RemoteAddr: remoteAddr,
		},
	}
}

// 本地段监听，接收来自本机浏览器的连接
func (local *LsLocal) Listen(didListen func(listenAddr net.Addr)) error {
	listener, err := net.ListenTCP("tcp", local.ListenAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	if didListen != nil {
		didListen(listener.Addr())
	}

	for {
		userConn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		// userConn被关闭的时候需要清除所有数据
		userConn.SetLinger(0)
		go local.handleConn(userConn)
	}
}

func (local *LsLocal) handleConn(userConn *net.TCPConn) {
	defer userConn.Close()

	proxyServer, err := local.DialRemote()
	if err != nil {
		log.Println(err)
		return
	}
	defer proxyServer.Close()
	// Conn被关闭的时候清除数据
	proxyServer.SetLinger(0)

	// 开始转发
	// proxyServer -> 读取数据 -> localUser
	go func() {
		err := local.DecodeCopy(userConn, proxyServer)
		if err != nil {
			// 在copy过程中，可能存在网络超时等待，error被return，只要有一个发生了错误就退出
			userConn.Close()
			proxyServer.Close()
		}
	}()

	// 从 localUser 发送数据发送到 proxyServer，这里因为处在翻墙阶段出现网络错误的概率更大
	local.EncodeCopy(proxyServer, userConn)
}
