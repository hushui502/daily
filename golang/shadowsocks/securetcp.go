package shadowsocks

import (
	"errors"
	"fmt"
	"io"
	"net"
)

const (
	BufSize = 1024
)

type SecureSocket struct {
	Cipher     *Cipher      // 密钥部分
	ListenAddr *net.TCPAddr // 监听的addr
	RemoteAddr *net.TCPAddr // 远称的addr
}

// 将输入流中的加密过的数据，解密后放在bs中
func (secureSocket *SecureSocket) DecodeRead(conn *net.TCPConn, bs []byte) (n int, err error) {
	n, err = conn.Read(bs)
	if err != nil {
		return
	}

	secureSocket.Cipher.decode(bs[:n])
	return
}

// 将放在bs中的未加密数据，加密后写入到输出流
func (secureSocket *SecureSocket) EncodeWrite(conn *net.TCPConn, bs []byte) (n int, err error) {
	secureSocket.Cipher.encode(bs)
	return conn.Write(bs)
}

// 从src中读原数据并加密写入dst，一直到src没有数据可以再read
func (secureSocket *SecureSocket) EncodeCopy(dst *net.TCPConn, src *net.TCPConn) error {
	buf := make([]byte, BufSize)
	for {
		// errRead这种写法代码阅读似乎更好一点，建议在复杂逻辑并且非常多err中尝试一下
		readCount, errRead := src.Read(buf)
		if errRead != nil {
			if errRead != io.EOF {
				return errRead
			} else {
				return nil
			}
		}
		if readCount > 0 {
			writeCount, errWrite := secureSocket.EncodeWrite(dst, buf[0:readCount])
			if errWrite != nil {
				return errWrite
			}
			if writeCount != readCount {
				return io.ErrShortWrite
			}
		}
	}
}

// 从src中读取加密后的数据解密后写入到dst，直到src中没有数据可以再读取
func (secureSocket *SecureSocket) DecodeCopy(dst *net.TCPConn, src *net.TCPConn) error {
	buf := make([]byte, BufSize)
	for {
		readCount, errRead := secureSocket.DecodeRead(src, buf)
		if errRead != nil {
			if errRead != io.EOF {
				return errRead
			} else {
				return nil
			}
		}
		if readCount > 0 {
			writeCount, errWrite := dst.Write(buf[0:readCount])
			if errWrite != nil {
				return errWrite
			}
			if writeCount != readCount {
				return io.ErrShortWrite
			}
		}
	}
}

// 和远称socket建立连接，之后他们之间的数据会传输加密
func (secureSocket *SecureSocket) DialRemote() (*net.TCPConn, error) {
	remoteConn, err := net.DialTCP("tcp", nil, secureSocket.RemoteAddr)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("连接到远程服务器 %s 失败:%s", secureSocket.RemoteAddr, err))
	}
	return remoteConn, nil
}
