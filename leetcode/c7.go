package leetcode


// 两数之和主要就是考察target-val这个思路
// 还有避免暴力的双循环,可以考虑map
func twoSum(nums []int, target int) []int {
	for i, v := range nums {
		for k := i+1; k < len(nums); k++ {
			if target-v == nums[k] {
				return []int{i, k}
			}
		}
	}
	return []int{}
}

func twoSum2(nums []int, target int) []int {
	m := make(map[int]int)
	for i, v := range nums {
		m[v] = i
	}

	for i, v := range nums {
		if val, ok := m[target-v]; ok {
			return []int{i, val}
		}
	}

	return []int{}
}
