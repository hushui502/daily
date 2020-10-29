package main

import (
	"testing"
	"time"
)

func BenchmarkLoad(b *testing.B) {
	for i := 0; i < b.N; i++ {

	}
}

func TestLoad(t *testing.T) {
	t.Parallel()
	time.Sleep(time.Second * 2)
}

func TestLoad2(t *testing.T) {
	t.Parallel()
	time.Sleep(time.Second * 3)
}
