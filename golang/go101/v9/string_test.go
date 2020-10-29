package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

const v = "Measure the elapsed time between sending a data octet with a?"

func BenchmarkStringJoin(b *testing.B) {
	var s string
	for i := 0; i < b.N; i++ {
		s = strings.Join([]string{s, "[", v, "]"}, "")
	}
}

func BenchmarkStringAdd(b *testing.B) {
	var s string
	for i := 0; i < b.N; i++ {
		s = s + "[" + v + "]"
	}
}

func BenchmarkSprintf(b *testing.B) {
	var s string
	for i := 0; i < b.N; i++ {
		s = fmt.Sprintf("%s[%s]", s, v)
	}
}

func BenchmarkBuffer(b *testing.B) {
	var buf bytes.Buffer
	for i := 0; i < b.N; i++ {
		buf.WriteString("[")
		buf.WriteString(v)
		buf.WriteString("]")
	}
}
