package leetcode

func twoSum(nums []int, target int) []int {
	if len(nums) == 0 {
		return []int{}
	}

	m := make(map[int]int)

	for i, v := range nums {
		if _, ok := m[target-v]; ok {
			return []int{m[target-v], i}
		}
		m[v] = i
	}

	return []int{}
}

func longestPalindrome(s string) string {
	if len(s) == 0 {
		return ""
	}

	start, end := 0, 0

	for i := 0; i < len(s); i++ {
		left1, right1 := expandAroundCenter(s, i, i)
		left2, right2 := expandAroundCenter(s, i, i+1)

		if right1-left1 > end-start {
			start, end = left1, right1
		}
		if right2-left2 > end-start {
			start, end = left2, right2
		}
	}

	return s[start : end+1]
}

func expandAroundCenter(s string, left, right int) (int, int) {
	for left >= 0 && right < len(s) && s[left] == s[right] {
		left--
		right++
	}

	return left + 1, right - 1
}

func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
	nums3 := append(nums1, nums2...)
	quickSortHelper(nums3, 0, len(nums3)-1)

	if len(nums3)%2 == 0 {
		return float64(nums3[len(nums3)/2]+nums3[len(nums3)/2-1]) / 2
	} else {
		return float64(nums3[len(nums3)/2])
	}
}
