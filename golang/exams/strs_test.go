package main

import (
	"strings"
	"sync"
	"testing"
)


func Test_isUniqueString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "str1",
			args: args{s: "abcdef"},
			want: true,
		},
		{
			name: "str2",
			args: args{s: "abcdefa"},
			want: false,
		},
	}
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isUniqueString(tt.args.s); got != tt.want {
				t.Errorf("isUniqueString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isUniqueString2(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "str1",
			args: args{s: "abcdef"},
			want: true,
		},
		{
			name: "str2",
			args: args{s: "abcdefa"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isUniqueString2(tt.args.s); got != tt.want {
				t.Errorf("isUniqueString2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_printA(t *testing.T) {
	type args struct {
		number chan bool
		letter chan bool
		wg     *sync.WaitGroup
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_reverseString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			name: "str1",
			args: args{s: "abcd"},
			want: "dcba",
			want1: true,
		},
		{
			name: "str2",
			args: args{s: "qwer"},
			want: "rewq",
			want1: true,
		},
		{
			name: "str3",
			args: args{s: strings.Repeat("abcde", 10000)},
			want: strings.Repeat("abcde", 10000),
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := reverseString(tt.args.s)
			if got != tt.want {
				t.Errorf("reverseString() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("reverseString() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_isRegroup(t *testing.T) {
	type args struct {
		s1 string
		s2 string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "s1==s2",
			args: args{
				s1: "abc",
				s2: "cba",
			},
			want: true,
		},
		{
			name: "s1!=s2",
			args: args{
				s1: "abe",
				s2: "cba",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRegroup(tt.args.s1, tt.args.s2); got != tt.want {
				t.Errorf("isRegroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_replaceBlank(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			name: "str1",
			args: args{s: "ab c"},
			want: "ab%20c",
			want1: true,
		},
		{
			name: "str2",
			args: args{s: "ab23c"},
			want: "ab23c",
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := replaceBlank(tt.args.s)
			if got != tt.want {
				t.Errorf("replaceBlank() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("replaceBlank() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}