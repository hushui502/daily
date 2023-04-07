package leetcode

import (
	"container/list"
	"path/filepath"
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
			name: "case1",
			args: args{
				nums:   []int{1, 2, 3, 2},
				target: 4,
			},
			want: []int{0, 2},
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

func TestlongestPalindrome(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case1",
			args: args{
				s: "babad",
			},
			want: "bab",
		},
		{
			name: "case2",
			args: args{
				s: "cbbd",
			},
			want: "bb",
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

func TestFilePath(t *testing.T) {
	path, err := filepath.Abs("solution_01_test.go")
	if err != nil {
		t.Error(err)
	}
	t.Log(path)
}

func Test_runningSum(t *testing.T) {
	type args struct {
		nums []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "case1",
			args: args{
				nums: []int{1, 2, 3, 4},
			},
			want: []int{1, 3, 6, 10},
		},
		{
			name: "case2",
			args: args{
				nums: []int{1, 1, 1, 1, 1},
			},
			want: []int{1, 2, 3, 4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := runningSum(tt.args.nums); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("runningSum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reverseList(t *testing.T) {
	type args struct {
		head *ListNode
	}
	tests := []struct {
		name string
		args args
		want *ListNode
	}{
		{
			name: "case1",
			args: args{
				head: &ListNode{
					Val: 1,
					Next: &ListNode{
						Val: 2,
						Next: &ListNode{
							Val: 3,
							Next: &ListNode{
								Val: 4,
							},
						},
					},
				},
			},
			want: &ListNode{
				Val: 4,
				Next: &ListNode{
					Val: 3,
					Next: &ListNode{
						Val: 2,
						Next: &ListNode{
							Val: 1,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reverseList(tt.args.head); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("reverseList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_zigzagLevelOrder(t *testing.T) {
	type args struct {
		root *TreeNode
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{
			name: "case1",
			args: args{
				root: &TreeNode{
					Val: 3,
					Left: &TreeNode{
						Val: 9,
					},
					Right: &TreeNode{
						Val: 20,
						Left: &TreeNode{
							Val: 15,
						},
						Right: &TreeNode{
							Val: 7,
						},
					},
				},
			},
			want: [][]int{
				{3},
				{20, 9},
				{15, 7},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := zigzagLevelOrder(tt.args.root); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("zigzagLevelOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func Test_lowestCommonAncestor(t *testing.T) {
//	type args struct {
//		root *TreeNode
//		p    *TreeNode
//		q    *TreeNode
//	}
//	tests := []struct {
//		name string
//		args args
//		want *TreeNode
//	}{
//		{
//			name: "case1",
//			args: args{
//				root: &TreeNode{
//					Val: 3,
//					Left: &TreeNode{
//						Val: 5,
//						Left: &TreeNode{
//							Val: 6,
//						},
//						Right: &TreeNode{
//							Val: 2,
//							Left: &TreeNode{
//								Val: 7,
//							},
//						},
//					},
//					Right: &TreeNode{
//						Val: 1,
//						Left: &TreeNode{
//							Val: 0,
//						},
//						Right: &TreeNode{
//							Val: 8,
//						},
//					},
//				},
//				p: &TreeNode{
//					Val: 5,
//					Left: &TreeNode{
//						Val: 6,
//					},
//					Right: &TreeNode{
//						Val: 2,
//						Left: &TreeNode{
//							Val: 7,
//						},
//					},
//				},
//				q: &TreeNode{
//					Val: 1,
//					Left: &TreeNode{
//						Val: 0,
//					},
//					Right: &TreeNode{
//						Val: 8,
//					},
//				},
//			},
//			want: &TreeNode{
//				Val: 3,
//				Left: &TreeNode{
//					Val: 5,
//					Left: &TreeNode{
//						Val: 6,
//					},
//					Right: &TreeNode{
//						Val: 2,
//						Left: &TreeNode{
//							Val: 7,
//						},
//					},
//				},
//				Right: &TreeNode{
//					Val: 1,
//					Left: &TreeNode{
//						Val: 0,
//					},
//					Right: &TreeNode{
//						Val: 8,
//					},
//				},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := lowestCommonAncestor(tt.args.root, tt.args.p, tt.args.q); got != tt.want {
//				t.Errorf("lowestCommonAncestor() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func Test_numIslands(t *testing.T) {
	type args struct {
		grid [][]byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "case1",
			args: args{
				grid: [][]byte{
					{'1', '1', '1', '1', '0'},
					{'1', '1', '0', '1', '0'},
					{'1', '1', '0', '0', '0'},
					{'0', '0', '0', '0', '0'},
				},
			},
			want: 1,
		},
		{
			name: "case2",
			args: args{
				grid: [][]byte{
					{'1', '1', '0', '0', '0'},
					{'1', '1', '0', '0', '0'},
					{'0', '0', '1', '0', '0'},
					{'0', '0', '0', '1', '1'},
				},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := numIslands(tt.args.grid); got != tt.want {
				t.Errorf("numIslands() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_merge(t *testing.T) {
	type args struct {
		nums1 []int
		m     int
		nums2 []int
		n     int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "case1",
			args: args{
				nums1: []int{1, 2, 3, 0, 0, 0},
				m:     3,
				nums2: []int{2, 5, 6},
				n:     3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			merge(tt.args.nums1, tt.args.m, tt.args.nums2, tt.args.n)
		})
	}
}

func Test_maxProfit(t *testing.T) {
	type args struct {
		prices []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "case1",
			args: args{
				prices: []int{7, 1, 5, 3, 6, 4},
			},
			want: 5,
		},
		{
			name: "case2",
			args: args{
				prices: []int{7, 6, 4, 3, 1},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxProfit(tt.args.prices); got != tt.want {
				t.Errorf("maxProfit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConstructor(t *testing.T) {
	type args struct {
		capacity int
	}
	tests := []struct {
		name string
		args args
		want LRUCache
	}{
		{
			name: "case1",
			args: args{
				capacity: 2,
			},
			want: LRUCache{
				capacity: 2,
				cache:    make(map[int]*list.Element),
				ll:       list.New(),
			},
		},
		{
			name: "case2",
			args: args{
				capacity: 3,
			},
			want: LRUCache{
				capacity: 3,
				cache:    make(map[int]*list.Element),
				ll:       list.New(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Constructor(tt.args.capacity); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Constructor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLRUCache_Get(t *testing.T) {
	type fields struct {
		capacity int
		cache    map[int]*list.Element
		ll       *list.List
	}
	type args struct {
		key int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "case1",
			fields: fields{
				capacity: 2,
				cache:    make(map[int]*list.Element),
				ll:       list.New(),
			},
			args: args{
				key: 1,
			},
			want: -1,
		},
		{
			name: "case2",
			fields: fields{
				capacity: 2,
				cache:    make(map[int]*list.Element),
				ll:       list.New(),
			},
			args: args{
				key: 2,
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &LRUCache{
				capacity: tt.fields.capacity,
				cache:    tt.fields.cache,
				ll:       tt.fields.ll,
			}
			if got := this.Get(tt.args.key); got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLRUCache_Put(t *testing.T) {
	type fields struct {
		capacity int
		cache    map[int]*list.Element
		ll       *list.List
	}
	type args struct {
		key   int
		value int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "case1",
			fields: fields{
				capacity: 2,
				cache:    make(map[int]*list.Element),
				ll:       list.New(),
			},
			args: args{
				key:   1,
				value: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &LRUCache{
				capacity: tt.fields.capacity,
				cache:    tt.fields.cache,
				ll:       tt.fields.ll,
			}
			this.Put(tt.args.key, tt.args.value)
		})
	}
}

func Test_findKthLargest(t *testing.T) {
	type args struct {
		nums []int
		k    int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "case1",
			args: args{
				nums: []int{3, 2, 1, 5, 6, 4},
				k:    2,
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findKthLargest(tt.args.nums, tt.args.k); got != tt.want {
				t.Errorf("findKthLargest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hasCycle(t *testing.T) {
	type args struct {
		head *ListNode
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "case1",
			args: args{
				head: &ListNode{
					Val: 3,
					Next: &ListNode{
						Val: 2,
						Next: &ListNode{
							Val: 0,
							Next: &ListNode{
								Val:  -4,
								Next: nil,
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "case2",
			args: args{
				head: &ListNode{
					Val: 1,
					Next: &ListNode{
						Val: 2,
						Next: &ListNode{
							Val:  3,
							Next: nil,
						},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasCycle(tt.args.head); got != tt.want {
				t.Errorf("hasCycle() = %v, want %v", got, tt.want)
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
			name: "case1",
			args: args{
				s: "()",
			},
			want: true,
		},
		{
			name: "case2",
			args: args{
				s: "()[]{}",
			},
			want: true,
		},
		{
			name: "case3",
			args: args{
				s: "(]",
			},
			want: false,
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

func Test_search(t *testing.T) {
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
			name: "case1",
			args: args{
				nums:   []int{-1, 0, 3, 5, 9, 12},
				target: 9,
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := search(tt.args.nums, tt.args.target); got != tt.want {
				t.Errorf("search() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_maxSubArray(t *testing.T) {
	type args struct {
		nums []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "case1",
			args: args{
				nums: []int{-2, 1, -3, 4, -1, 2, 1, -5, 4},
			},
			want: 6,
		},
		{
			name: "case2",
			args: args{
				nums: []int{1},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxSubArray(tt.args.nums); got != tt.want {
				t.Errorf("maxSubArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mergeTwoLists(t *testing.T) {
	type args struct {
		list1 *ListNode
		list2 *ListNode
	}
	tests := []struct {
		name string
		args args
		want *ListNode
	}{
		{
			name: "case1",
			args: args{
				list1: &ListNode{
					Val: 1,
					Next: &ListNode{
						Val: 2,
						Next: &ListNode{
							Val:  4,
							Next: nil,
						},
					},
				},
				list2: &ListNode{
					Val: 1,
					Next: &ListNode{
						Val: 3,
						Next: &ListNode{
							Val:  4,
							Next: nil,
						},
					},
				},
			},
			want: &ListNode{
				Val: 1,
				Next: &ListNode{
					Val: 1,
					Next: &ListNode{
						Val: 2,
						Next: &ListNode{
							Val: 3,
							Next: &ListNode{
								Val: 4,
								Next: &ListNode{
									Val:  4,
									Next: nil,
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeTwoLists(tt.args.list1, tt.args.list2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeTwoLists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortArray(t *testing.T) {
	type args struct {
		nums []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "case1",
			args: args{
				nums: []int{5, 2, 3, 1},
			},
			want: []int{1, 2, 3, 5},
		},
		{
			name: "case2",
			args: args{
				nums: []int{5, 1, 1, 2, 0, 0},
			},
			want: []int{0, 0, 1, 1, 2, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sortArray(tt.args.nums); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortArray() = %v, want %v", got, tt.want)
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
			name: "case1",
			args: args{
				nums: []int{-1, 0, 1, 2, -1, -4},
			},
			want: [][]int{
				{-1, -1, 2},
				{0, 1, -1},
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

//func Test_getIntersectionNode(t *testing.T) {
//	type args struct {
//		headA *ListNode
//		headB *ListNode
//	}
//	tests := []struct {
//		name string
//		args args
//		want *ListNode
//	}{
//		{
//			name: "case1",
//			args: args{
//				headA: &ListNode{
//					Val: 4,
//					Next: &ListNode{
//						Val: 1,
//						Next: &ListNode{
//							Val: 8,
//							Next: &ListNode{
//								Val: 4,
//								Next: &ListNode{
//									Val: 5,
//									Next: &ListNode{
//										Val:  6,
//										Next: nil,
//									},
//								},
//							},
//						},
//					},
//				},
//				headB: &ListNode{
//					Val: 8,
//					Next: &ListNode{
//						Val: 4,
//						Next: &ListNode{
//							Val:  5,
//							Next: nil,
//						},
//					},
//				},
//			},
//			want: &ListNode{
//				Val: 8,
//				Next: &ListNode{
//					Val: 4,
//					Next: &ListNode{
//						Val:  5,
//						Next: nil,
//					},
//				},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := getIntersectionNode(tt.args.headA, tt.args.headB); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("getIntersectionNode() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func Test_spiralOrder(t *testing.T) {
	type args struct {
		matrix [][]int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "case1",
			args: args{
				matrix: [][]int{
					{1, 2, 3},
					{4, 5, 6},
					{7, 8, 9},
				},
			},
			want: []int{1, 2, 3, 6, 9, 8, 7, 4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := spiralOrder(tt.args.matrix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("spiralOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}
