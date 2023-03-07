package leetcode

func lengthOfLongestSubstring(s string) int {
	if len(s) == 0 {
		return 0
	}

	m := make(map[byte]int)
	left, maxLen := 0, 0

	for i := 0; i < len(s); i++ {
		if _, ok := m[s[i]]; ok {
			left = max(left, m[s[i]]+1)
		}
		m[s[i]] = i
		maxLen = max(maxLen, i-left+1)
	}

	return maxLen
}
