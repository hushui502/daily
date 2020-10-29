package v5

import (
	"io"
	"os"
)

var fd *os.File

func Write(filename string) error {
	var err error
	_, err = fd.Write([]byte("a"))
	if err != nil {
		return err
	}
	_, err = fd.Write([]byte("b"))
	if err != nil {
		return err
	}
	return nil
}

type errWriter struct {
	w   io.Writer
	err error
}

func (ew *errWriter) write(buf []byte) {
	if ew.err != nil {
		return
	}
	_, ew.err = ew.w.Write(buf)
}

func WriteFile1(filename string) error {
	ew := &errWriter{w: fd}
	ew.write([]byte("a"))
	ew.write([]byte("b"))
	if ew.err != nil {
		return ew.err
	}
	return nil
}
