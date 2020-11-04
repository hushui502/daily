package main

func firstMissingPositive(nums []int) int {
	numsMap := make(map[int]int, len(nums))
	for _, v := range numsMap {
		numsMap[v] = v
	}

	for index := 1; index < len(nums)+1; index++ {
		if _, ok := numsMap[index]; !ok {
			return index
		}
	}

	return len(nums) + 1
}