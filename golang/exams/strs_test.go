package main

import (
	"reflect"
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
			name:  "str1",
			args:  args{s: "abcd"},
			want:  "dcba",
			want1: true,
		},
		{
			name:  "str2",
			args:  args{s: "qwer"},
			want:  "rewq",
			want1: true,
		},
		{
			name:  "str3",
			args:  args{s: strings.Repeat("abcde", 10000)},
			want:  strings.Repeat("abcde", 10000),
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
			name:  "str1",
			args:  args{s: "ab c"},
			want:  "ab%20c",
			want1: true,
		},
		{
			name:  "str2",
			args:  args{s: "ab23c"},
			want:  "ab23c",
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

func Test_firstUniqueChar(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "str1",
			args: args{s: "leetcode"},
			want: 0,
		},
		{
			name: "str2",
			args: args{s: "loveleetcode"},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := firstUniqueChar(tt.args.s); got != tt.want {
				t.Errorf("firstUniqueChar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isPalindrome(t *testing.T) {
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
			args: args{s: "abcba"},
			want: true,
		},
		{
			name: "str2",
			args: args{s: "abb"},
			want: false,
		},
		{
			name: "str3",
			args: args{s: "ab1ba"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPalindrome(tt.args.s); got != tt.want {
				t.Errorf("isPalindrome() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_maxSlideWindow(t *testing.T) {
	type args struct {
		nums []int
		k    int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "win1",
			args: args{
				nums: []int{1, 3, -1, -3, 5, 3, 6, 7},
				k:    3,
			},
			want: []int{3, 3, 5, 5, 6, 7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxSlideWindow(tt.args.nums, tt.args.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("maxSlideWindow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bubbleSort(t *testing.T) {
	type args struct {
		arr []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "sort1",
			args: args{arr: []int{3, 1, 2, 4, 6}},
			want: []int{6, 4, 3, 2, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bubbleSort(tt.args.arr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bubbleSort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_selectSort(t *testing.T) {
	type args struct {
		arr []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "sort1",
			args: args{arr: []int{3, 1, 2, 4, 6}},
			want: []int{1, 2, 3, 4, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := selectSort(tt.args.arr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("selectSort() = %v, want %v", got, tt.want)
			}
		})
	}
}
