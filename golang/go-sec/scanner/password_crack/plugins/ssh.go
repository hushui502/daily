package plugins

import (
	"fmt"
	"go-sec/scanner/password_crack/models"
	"go-sec/scanner/password_crack/vars"
	"golang.org/x/crypto/ssh"
	"net"
)

func ScanSsh(s models.Service) (result models.ScanResult, err error) {
	result.Service = s
	config := &ssh.ClientConfig{
		User: s.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		Timeout: vars.Timeout,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", s.Ip, s.Port), config)
	if err != nil {
		return result, err
	}
	defer func() {
		if client != nil {
			client.Close()
		}
	}()

	session, err := client.NewSession()
	if err != nil {
		return result, err
	}
	defer func() {
		if session != nil {
			session.Close()
		}
	}()

	err = session.Run("echo 333")
	if err != nil {
		return result, err
	}

	result.Result = true

	return result, err
}
