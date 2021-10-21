package main

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParseAddr(t *testing.T) {
	ip, port := parseAddr("2E0010AC:E898")

	assert.Equal(t, ip, "172.16.0.46")
	assert.Equal(t, port, "59544")
}

func TestHandleFile(t *testing.T) {
	for _, fdn := range []string{"0", "1", "2", "3", "4"} {
		name := "/proc/994998/fd/" + fdn

		kk, err := os.Readlink(name)
		fmt.Println(kk, err)

		//只是针对 unix socket
		fileinfo, err := os.Lstat(name)
		assert.Equal(t, nil, err)
		name = "/tmp/mongodb-27017.sock"
		fmt.Println("socket:  ", fileinfo.Mode()&os.ModeSocket == os.ModeSocket)
		fmt.Println("symlink:  ", fileinfo.Mode()&os.ModeSymlink == os.ModeSymlink)

		if fileinfo.Mode()&os.ModeSocket == os.ModeSocket {
			fmt.Println("this is socket")
		}
		if fileinfo.Mode() == os.ModeSocket {
			fmt.Println("this is socket")
		}
		fmt.Println(fileinfo.Mode().String())

		fmt.Println("socket333:  ", fileinfo.Mode())
		fmt.Println("socket444:  ", KindFromFileInfo(fileinfo) == SOCKET)
		fmt.Println("socket555:  ", fileinfo.Mode()&os.ModeSocket)
	}
}

type Kind int

const (
	DIR Kind = iota
	LINK
	PIPE
	SOCKET
	DEV
	FILE
)

func KindFromFileInfo(fileInfo os.FileInfo) Kind {
	if fileInfo.IsDir() {
		return DIR
	}

	if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
		return LINK
	}

	if fileInfo.Mode()&os.ModeNamedPipe == os.ModeNamedPipe {
		return PIPE
	}

	if fileInfo.Mode()&os.ModeSocket == os.ModeSocket {
		return SOCKET
	}

	if fileInfo.Mode()&os.ModeDevice == os.ModeDevice {
		return DEV
	}

	return FILE
}

func (self Kind) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.String())
}

func (self *Kind) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case "DIR":
		*self = DIR
	case "LINK":
		*self = LINK
	case "PIPE":
		*self = PIPE
	case "SOCKET":
		*self = SOCKET
	case "DEV":
		*self = DEV
	case "FILE":
		*self = FILE
	default:
		return fmt.Errorf("invalid Kind: '%s'", s)
	}
	return nil
}

func (self Kind) String() string {
	names := []string{
		"DIR", "LINK", "PIPE", "SOCKET", "DEV", "FILE",
	}

	// FIXME: bound check?
	return names[self]
}
