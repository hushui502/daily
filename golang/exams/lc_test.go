package main

import (
	"reflect"
	"testing"
)

func Test_twoSum(t *testing.T) {
	type args struct {
		nums   []int
		target int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "test1",
			args: args{
				nums:   []int{1, 2, 3, 4},
				target: 3,
			},
			want: []int{1, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := twoSum(tt.args.nums, tt.args.target); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("twoSum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_zTraversalTree(t *testing.T) {
	type args struct {
		root *TreeNode
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := zTraversalTree(tt.args.root); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("zTraversalTree() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_searchRotateSortedArray(t *testing.T) {
	type args struct {
		nums   []int
		target int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "sort1",
			args: args{
				nums:   []int{4, 5, 6, 7, 0, 1, 2},
				target: 0,
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := searchRotateSortedArray(tt.args.nums, tt.args.target); got != tt.want {
				t.Errorf("searchRotateSortedArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lengthOfLongestSubstring(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test_lengthOfLongestSubstring",
			args: args{s: "abcabcbb"},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lengthOfLongestSubstring(tt.args.s); got != tt.want {
				t.Errorf("lengthOfLongestSubstring() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_threeSum(t *testing.T) {
	type args struct {
		nums []int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{
			name: "threeSum",
			args: args{nums: []int{-1, 0, 1, 2, -1, -4}},
			want: [][]int{
				{-1, 0, 1},
				{-1, -1, 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := threeSum(tt.args.nums); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("threeSum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_longestPalindrome(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "longestPalindrome",
			args: args{s: "babad"},
			want: "bab",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := longestPalindrome(tt.args.s); got != tt.want {
				t.Errorf("longestPalindrome() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_maxArea(t *testing.T) {
	type args struct {
		nums []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "maxArea",
			args: args{nums: []int{1, 8, 6, 2, 5, 4, 8, 3, 7}},
			want: 49,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxArea(tt.args.nums); got != tt.want {
				t.Errorf("maxArea() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findCombinations(t *testing.T) {
	type args struct {
		digits string
		index  int
		s      string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "combinations",
			args: args{
				digits: "23",
				index:  0,
				s:      "",
			},
			want: []string{"ad", "ae", "af", "bd", "be", "bf", "cd", "ce", "cf"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := letterCombinations(tt.args.digits); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("letters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removeNthElementFromEnd(t *testing.T) {
	type args struct {
		head *ListNode
		n    int
	}
	head := &ListNode{Val: 1,
		Next: &ListNode{Val: 2,
			Next: &ListNode{Val: 3,
				Next: &ListNode{Val: 4,
					Next: &ListNode{Val: 5,
						Next: nil}}}}}
	res := &ListNode{Val: 1,
		Next: &ListNode{Val: 2,
			Next: &ListNode{Val: 3,
				Next: &ListNode{Val: 5,
					Next: nil}}}}
	tests := []struct {
		name string
		args args
		want *ListNode
	}{
		{
			name: "Test_removeNthElementFromEnd",
			args: args{
				head: head,
				n:    2,
			},
			want: res,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeNthElementFromEnd(tt.args.head, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeNthElementFromEnd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isValid(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "isValid",
			args: args{s: "()"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValid(tt.args.s); got != tt.want {
				t.Errorf("isValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mergeTwoLists(t *testing.T) {
	type args struct {
		l1 *ListNode
		l2 *ListNode
	}
	l1 := &ListNode{Val: 1,
		Next: &ListNode{Val: 2,
			Next: &ListNode{Val: 3}}}
	l2 := &ListNode{Val: 4,
		Next: &ListNode{Val: 5,
			Next: nil}}

	res := &ListNode{Val: 1,
		Next: &ListNode{Val: 2,
			Next: &ListNode{Val: 3,
				Next: &ListNode{Val: 4,
					Next: &ListNode{Val: 5,
						Next: nil}}}}}

	tests := []struct {
		name string
		args args
		want *ListNode
	}{
		{
			name: "mergeList",
			args: args{
				l1: l1,
				l2: l2,
			},
			want: res,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeTwoLists(tt.args.l1, tt.args.l2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeTwoLists2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mergeTwoLists2(t *testing.T) {
	type args struct {
		l1 *ListNode
		l2 *ListNode
	}
	l1 := &ListNode{Val: 1,
		Next: &ListNode{Val: 2,
			Next: &ListNode{Val: 3}}}
	l2 := &ListNode{Val: 4,
		Next: &ListNode{Val: 5,
			Next: nil}}

	res := &ListNode{Val: 1,
		Next: &ListNode{Val: 2,
			Next: &ListNode{Val: 3,
				Next: &ListNode{Val: 4,
					Next: &ListNode{Val: 5,
						Next: nil}}}}}

	tests := []struct {
		name string
		args args
		want *ListNode
	}{
		{
			name: "mergeList",
			args: args{
				l1: l1,
				l2: l2,
			},
			want: res,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeTwoLists2(tt.args.l1, tt.args.l2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeTwoLists2() = %v, want %v", got, tt.want)
			}
		})
	}
}
