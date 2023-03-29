package main

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
