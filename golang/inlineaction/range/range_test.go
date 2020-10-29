package main

import "testing"

const size = 1000000

type SomeStruct struct {
	ID0 int64
	ID1 int64
	ID2 int64
	ID3 int64
	ID4 int64
	ID5 int64
	ID6 int64
	ID7 int64
	ID8 int64
	//ID9 int64
}

func BenchmarkForVar(b *testing.B) {
	slice := make([]SomeStruct, size)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, s := range slice { // index and value
			_ = s
		}
	}
}
func BenchmarkForCounter(b *testing.B) {
	slice := make([]SomeStruct, size)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i := range slice { // only use the index
			s := slice[i]
			_ = s
		}
	}
}

func BenchmarkFor(b *testing.B) {
	slice := make([]SomeStruct, size)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		//slice := make([]SomeStruct, size)
		for i := 0; i < len(slice); i++ {
			s := slice[i]
			_ = s
		}
	}
}
