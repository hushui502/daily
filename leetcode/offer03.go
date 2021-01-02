package leetcode

import "sort"

func findRepeatNumber1(nums []int) int {
	m := make(map[int]int)

	for _, v := range nums {
		if _, ok := m[v]; ok {
			return v
		} else {
			m[v] = 1
		}
	}

	return -1
}

func findRepeatNumber2(nums []int) int {
	sort.Ints(nums)
	numSize := len(nums)
	for i := 0; i < numSize-1; i++ {
		if nums[i] == nums[i+1] {
			return nums[i]
		}
	}

	return -1
}

func findRepeatNumber3(nums []int) int {
	for i := 0; i < len(nums); i++ {
		for i != nums[i] {
			if nums[i] == nums[nums[i]] {
				return nums[i]
			}
			nums[i], nums[nums[i]] = nums[nums[i]], nums[i]
		}
	}

	return -1
}