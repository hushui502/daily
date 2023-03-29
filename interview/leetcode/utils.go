package main

type ListNode struct {
	Val  int
	Next *ListNode
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func min(a, b int) int {
	if a > b {
		return b
	}

	return a
}

func abs(a int) int {
	if a < 0 {
		return -a
	}

	return a
}

func quickSort(nums []int) {
	if len(nums) == 0 {
		return
	}

	quickSortHelper(nums, 0, len(nums)-1)
}

func quickSortHelper(nums []int, left, right int) {
	if left >= right {
		return
	}

	pivot := partition(nums, left, right)
	quickSortHelper(nums, left, pivot-1)
	quickSortHelper(nums, pivot+1, right)
}

func partition(nums []int, left, right int) int {
	pivot := nums[right]
	i := left
	for j := left; j < right; j++ {
		if nums[j] < pivot {
			nums[i], nums[j] = nums[j], nums[i]
			i++
		}
	}
	nums[i], nums[right] = nums[right], nums[i]

	return i
}
