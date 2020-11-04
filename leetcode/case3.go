package main

func lengthOfLongestSubstring(s string) int {
	if len(s) == 0 {
		return 0
	}

	var bitSet [256]bool
	result, left, right := 0, 0, 0
	for left < len(s) {
		if bitSet[s[right]] {
			bitSet[s[left]] = false
			left++
		} else {
			bitSet[s[right]] = true
			right++
		}
		if result < right - left {
			result = right - left
		}
		if result+left >= len(s) || right >= len(s) {
			break
		}
	}
	return result
}

func lengthOfLongestSubstring2(s string) int {
	if len(s) == 0 {
		return 0
	}

	var freq [256]int
	result, left, right := 0, 0, -1

	for left < len(s) {
		if right+1 < len(s) && freq[s[right+1]-'a'] == 0 {
			freq[s[right+1]-'a']++
			right++
		} else {
			freq[s[left]-'a']++
			left++
		}
		result = max(result, right-left+1)
	}
	return result
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
