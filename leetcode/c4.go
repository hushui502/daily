package leetcode

func reverse(arr []int) {
	for i := 0; i < len(arr)/2; i++ {
		arr[i], arr[len(arr)-1-i] = arr[len(arr)-1-i], arr[i]
	}
}

//输入: [1,2,3,4,5,6,7] 和 k = 3
//输出: [5,6,7,1,2,3,4]

func rotate(nums []int, k int) {
	// 首先反转整体
	reverse(nums)
	// 反转k之前
	reverse(nums[:k%len(nums)])
	// 反转k之后
	reverse(nums[k%len(nums):])
}
