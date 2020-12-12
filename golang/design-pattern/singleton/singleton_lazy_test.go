package singleton

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetLazySingelton(t *testing.T) {
	assert.Equal(t, GetLazySingelton(), GetLazySingelton())
}

func BenchmarkGetLazySingelton(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if GetLazySingelton() != GetLazySingelton() {
				b.Error("test failed")
			}
		}
	})
}
