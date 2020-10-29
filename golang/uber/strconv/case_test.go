package main

import (
	"os"
	"testing"
)

func BenchmarkTest(b *testing.B) {
	w, _ := os.Open("uber/test")
	for i := 0; i < b.N; i++ {
		w.Write([]byte("Hello world"))
	}
}

func BenchmarkTest2(b *testing.B) {
	data := []byte("hello world")
	w, _ := os.Open("D:\\project\\go\\src\\awesomeProject2\\uber\\test")
	for i := 0; i < b.N; i++ {
		w.Write(data)
	}
}
