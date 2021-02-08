package core

import (
	"io"
	"log"
	"net"
	"sync"
)

const (
	BufSize = 1024
)

var bpool sync.Pool

func init() {
	bpool.New = func() interface{} {
		return make([]byte, BufSize)
	}
}

func bufferPoolGet() []byte {
	return bpool.Get().([]byte)
}

func bufferPoolPut(b []byte) {
	bpool.Put(b)
}

type SecureTCPConn struct {
	io.ReadWriteCloser
	Cipher     *Cipher      // 密钥部分
	//ListenAddr *net.TCPAddr // 监听的addr
	//RemoteAddr *net.TCPAddr // 远称的addr
}

// 将输入流中的加密过的数据，解密后放在bs中
func (secureSocket *SecureTCPConn) DecodeRead(bs []byte) (n int, err error) {
	n, err = secureSocket.Read(bs)
	if err != nil {
		return
	}

	secureSocket.Cipher.decode(bs[:n])
	return
}

// 将放在bs中的未加密数据，加密后写入到输出流
func (secureSocket *SecureTCPConn) EncodeWrite(bs []byte) (n int, err error) {
	secureSocket.Cipher.encode(bs)
	return secureSocket.Write(bs)
}

// 从src中读原数据并加密写入dst，一直到src没有数据可以再read
func (secureSocket *SecureTCPConn) EncodeCopy(dst io.ReadWriteCloser) error {
	buf := bufferPoolGet()
	defer bufferPoolPut(buf)
	for {
		// errRead这种写法代码阅读似乎更好一点，建议在复杂逻辑并且非常多err中尝试一下
		readCount, errRead := secureSocket.Read(buf)
		if errRead != nil {
			if errRead != io.EOF {
				return errRead
			} else {
				return nil
			}
		}
		if readCount > 0 {
			writeCount, errWrite := (&SecureTCPConn{
				ReadWriteCloser: dst,
				Cipher:          secureSocket.Cipher,
			}).EncodeWrite(buf[0:readCount])
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
func (secureSocket *SecureTCPConn) DecodeCopy(dst io.Writer) error {
	buf := bufferPoolGet()
	defer bufferPoolPut(buf)
	for {
		readCount, errRead := secureSocket.DecodeRead(buf)
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

func DialEncryptedTCP(raddr *net.TCPAddr, cipher *Cipher) (*SecureTCPConn, error) {
	remoteConn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		return nil, err
	}
	remoteConn.SetLinger(0)

	return &SecureTCPConn{
		ReadWriteCloser: remoteConn,
		Cipher:          cipher,
	}, nil
}

func ListenEncryptedTCP(laddr *net.TCPAddr, cipher *Cipher, handleConn func(localConn *SecureTCPConn), didListen func(listenAddr *net.TCPAddr)) error {
	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	if didListen != nil {
		// didListen 可能有阻塞操作
		go didListen(listener.Addr().(*net.TCPAddr))
	}

	for {
		localConn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		// localConn被关闭时直接清除所有数据 不管没有发送的数据
		localConn.SetLinger(0)
		go handleConn(&SecureTCPConn{
			ReadWriteCloser: localConn,
			Cipher:          cipher,
		})
	}
}