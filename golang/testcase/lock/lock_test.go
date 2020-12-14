package lock

import (
	"sync"
	"testing"
)

func benchmark(b *testing.B, rw RW, read, write int) {
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for k := 0; k < read*100; k++ {
			wg.Add(1)
			go func() {
				rw.Read()
				wg.Done()
			}()
		}
		for k := 0; k < write*100; k++ {
			wg.Add(1)
			go func() {
				rw.Write()
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

//				const cost = time.Microsecond
//BenchmarkReadMore-8                    1        1011109700 ns/op
//BenchmarkReadMoreRW-8                 10         104401360 ns/op
//BenchmarkWriteMore-8                   1        1015632800 ns/op
//BenchmarkWriteMoreRW-8                 2         923991550 ns/op
//BenchmarkEqual-8                       1        1031583800 ns/op
//BenchmarkEqualRW-8                     2         513871450 ns/op
// 				实际上两者如果在读写请求并行且次数几乎相同时，差距不大


func BenchmarkReadMore(b *testing.B)    { benchmark(b, &Lock{}, 9, 1) }
func BenchmarkReadMoreRW(b *testing.B)  { benchmark(b, &RWLock{}, 9, 1) }
func BenchmarkWriteMore(b *testing.B)   { benchmark(b, &Lock{}, 1, 9) }
func BenchmarkWriteMoreRW(b *testing.B) { benchmark(b, &RWLock{}, 1, 9) }
func BenchmarkEqual(b *testing.B)       { benchmark(b, &Lock{}, 5, 5) }
func BenchmarkEqualRW(b *testing.B)     { benchmark(b, &RWLock{}, 5, 5) }


