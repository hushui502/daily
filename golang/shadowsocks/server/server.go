package server

import (
	"encoding/binary"
	"net"
	"shadowsocks/core"
)

// 新建一个服务端
// 服务端的职责是:
// 1. 监听来自本地代理客户端的请求
// 2. 解密本地代理客户端请求的数据，解析 SOCKS5 协议，连接用户浏览器真正想要连接的远程服务器
// 3. 转发用户浏览器真正想要连接的远程服务器返回的数据的加密后的内容到本地代理客户端
type LsServer struct {
	Cipher *core.Cipher
	ListenAddr *net.TCPAddr
}

func NewLsServer(password string, listenAddr string) (*LsServer, error) {
	bsPassword, err := core.ParsePassword(password)
	if err != nil {
		return nil, err
	}
	structListenAddr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	return &LsServer{
		Cipher:     core.NewCipher(bsPassword),
		ListenAddr: structListenAddr,
	}, nil
}

func (lsServer *LsServer) Listen(didListen func(listenAddr *net.TCPAddr)) error {
	return core.ListenEncryptedTCP(lsServer.ListenAddr, lsServer.Cipher, lsServer.handleConn, didListen)
}

// 解析socks5协议
// https://www.ietf.org/rfc/rfc1928.txt
func (lsServer *LsServer) handleConn(localConn *core.SecureTCPConn) {
	defer localConn.Close()
	buf := make([]byte, 256)
	/**
	   The localConn connects to the dstServer, and sends a ver
	   identifier/method selection message:
		          +----+----------+----------+
		          |VER | NMETHODS | METHODS  |
		          +----+----------+----------+
		          | 1  |    1     | 1 to 255 |
		          +----+----------+----------+
	   The VER field is set to X'05' for this ver of the protocol.  The
	   NMETHODS field contains the number of method identifier octets that
	   appear in the METHODS field.
	*/
	// 第一个字段VER代表版本号，SOCKS5默认位0x05，其固定长度为1字节
	_, err := localConn.DecodeRead(buf)
	// 只支持版本5
	if err != nil || buf[0] != 0x05 {
		return
	}

	/**
	   The dstServer selects from one of the methods given in METHODS, and
	   sends a METHOD selection message:

		          +----+--------+
		          |VER | METHOD |
		          +----+--------+
		          | 1  |   1    |
		          +----+--------+
	*/
	// 不需要验证
	localConn.EncodeWrite([]byte{0x05, 0x06})

	/**
	  +----+-----+-------+------+----------+----------+
	  |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	  +----+-----+-------+------+----------+----------+
	  | 1  |  1  | X'00' |  1   | Variable |    2     |
	  +----+-----+-------+------+----------+----------+
	*/
	// 获取真正的远称服务地址
	n, err := localConn.DecodeRead(buf)
	// n 最短的长度为7 情况是ATY=3 DST.ADDR占用一个字节 值为0x0
	if err != nil || n < 7 {
		return
	}

	// CMD 代表客户端请求的类型，值长度为1个字节，有三种类型
	if buf[1] != 0x01 {
		// 目前只支持 CONNECT
		return
	}

	var dIP []byte
	// aType 代表了请求远称服务器地址的类型，值长度1个字节
	switch buf[3] {
	case 0x01:
		// IPV4
		dIP = buf[4:4+net.IPv4len]
	case 0x03:
		ipAddr, err := net.ResolveIPAddr("ip", string(buf[5:n-2]))
		if err != nil {
			return
		}
		dIP = ipAddr.IP
	case 0x04:
		dIP = buf[4:4+net.IPv6len]
	default:
		return
	}
	dPort := buf[n-2:]
	dstAddr := &net.TCPAddr{
		IP:   dIP,
		Port: int(binary.BigEndian.Uint16(dPort)),
	}

	// 连接真正的服务
	dstServer, err := net.DialTCP("tcp", nil, dstAddr)
	if err != nil {
		return
	} else {
		defer dstServer.Close()
		dstServer.SetLinger(0)
		// 响应客户端连接成功
		/**
		  +----+-----+-------+------+----------+----------+
		  |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
		  +----+-----+-------+------+----------+----------+
		  | 1  |  1  | X'00' |  1   | Variable |    2     |
		  +----+-----+-------+------+----------+----------+
		*/
		// 响应客户端连接成功
		localConn.EncodeWrite([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	}

	// 进行转发
	// localConn -> read content -> dstServer
	go func() {
		err := localConn.DecodeCopy(dstServer)
		if err != nil {
			localConn.Close()
			dstServer.Close()
		}
	}()
	// 从 dstServer 读取数据发送到 localUser，这里因为处在翻墙阶段出现网络错误的概率更大
	(&core.SecureTCPConn{
		Cipher:          localConn.Cipher,
		ReadWriteCloser: dstServer,
	}).EncodeCopy(localConn)
}

