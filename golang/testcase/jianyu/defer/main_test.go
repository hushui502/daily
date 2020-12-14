package main

import "testing"

func BenchmarkDoDefer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DoDefer("hufan", "hhh")
	}
}

func BenchmarkDoNotDefer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DoNotDefer("煎鱼", "https://github.com/EDDYCJY/blog")
	}
}
