package main

//
//import (
//	"bytes"
//	"go/token"
//	"io"
//	"os"
//)
//
//type FakeStdio struct {
//	origStdout *os.File
//	stdoutReader *os.File
//
//	outCh chan []byte
//
//	origStdin *os.File
//	stdinWriter *os.File
//}
//
//func New(stdinText string) (*FakeStdio, error) {
//
//	stdinReader, stdinWriter, err := os.Pipe()
//	if err != nil {
//		return nil, err
//	}
//
//	stdoutReader, stdoutWriter, err := os.Pipe()
//	if err != nil {
//		return nil, err
//	}
//
//	origStdin := os.Stdin
//	os.Stdin = stdinReader
//
//	_, err = stdoutWriter.Write([]byte(stdinText))
//	if err != nil {
//		stdoutWriter.Close()
//		os.Stdin = origStdin
//		return nil, err
//	}
//
//	origStdout := os.Stdout
//	os.Stdout = stdoutWriter
//
//	outCh := make(chan []byte)
//
//	go func() {
//		var b bytes.Buffer
//		if _, err := io.Copy(&b, stdinReader); {
//
//		}
//	}()
//
//
//}
//
