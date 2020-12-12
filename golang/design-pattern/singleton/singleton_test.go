package singleton

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSingleton(t *testing.T) {
	assert.Equal(t, GetSingleton(), GetSingleton())
}

func BenchmarkGetSingleton(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if GetSingleton() != GetSingleton() {
				b.Error("test failed")
			}
		}
	})
}
