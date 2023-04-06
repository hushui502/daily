package leetcode

func lengthOfLongestSubstring(s string) int {
	if len(s) == 0 {
		return 0
	}

	m := make(map[byte]int)
	left, right := 0, 0
	maxLen := 0

	for right < len(s) {
		if _, ok := m[s[right]]; ok {
			delete(m, s[left])
			left++
		} else {
			m[s[right]] = right
			right++
		}

		maxLen = max(maxLen, right-left)
	}

	return maxLen
}

func lengthOfLongestSubstring2(s string) int {
	if len(s) == 0 {
		return 0
	}

	bs := [256]byte{}
	left, right := 0, -1
	maxLen := 0

	for right+1 < len(s) {
		if right+1 < len(s) && bs[s[right+1]] == 0 {
			bs[s[right+1]] = 1
			right++
		} else {
			bs[s[left]] = 0
			left++
		}

		maxLen = max(maxLen, right-left+1)
	}

	return maxLen
}
