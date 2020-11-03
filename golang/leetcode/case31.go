package main

import "sort"

func combinationSum(candidates []int, target int) [][]int {
	if len(candidates) == 0 {
		return [][]int{}
	}
	c, res := []int{}, [][]int{}
	sort.Ints(candidates)
	findCombinationSum(candidates, target, 0, c, &res)
	return res
}

func findCombinationSum(nums []int, target, index int, c []int, res *[][]int) {
	if target <= 0 {
		if target == 0 {
			b := make([]int, len(c))
			copy(b, c)
			*res = append(*res, b)
		}
		return
	}

	for i := index; i < len(nums); i++ {
		if nums[i] > target {
			break
		}
		c = append(c, nums[i])
		findCombinationSum(nums, target-nums[i], i, c, res)
		c = c[:len(c)-1]
	}
}

// eacg member in candidates may only be used once in the combination
func findCombinationSum2(nums []int, target, index int, c []int, res *[][]int) {
	if target == 0 {
		b := make([]int, len(c))
		copy(b, c)
		*res = append(*res, b)
		return
	}

	for i := index; i < len(nums); i++ {
		if i > index && nums[i] == nums[i-1] {
			continue
		}
		if target >= nums[i] {
			c = append(c, nums[i])
			findCombinationSum2(nums, target-nums[i], i+1, c, res)
			c = c[:len(c)-1]
		}
	}
}













