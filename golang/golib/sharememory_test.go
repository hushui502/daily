package golib

import "testing"

func BenchmarkPad(b *testing.B) {
	s := pad{}
	//p := pdd{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.increase()
			//p.increase()
		}
	})
}

//BenchmarkPad-8   	18799878	        62.5 ns/op
