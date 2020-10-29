package main

import "testing"

func BenchmarkUse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		use()
	}
}
