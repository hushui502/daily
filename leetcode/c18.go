package leetcode

import "strings"

func reverseString(s []byte) {
	left := 0
	right := len(s) - 1
	for left <= right {
		// 借助Golang的feature
		s[left], s[right] = s[right], s[left]
		left++
		right--
	}
}

func firstUniqChar(s string) int {
	// 如果用map的话是无法保证顺序的
	charCount := [26]int{}
	for _, ch := range s {
		charCount[ch - 'a']++
	}
	for i, ch := range s {
		if charCount[ch] == 1 {
			return i
		}
	}

	return -1
}

func strStr(haystack string, needle string) int {
	l1 := len(haystack)
	l2 := len(needle)

	if l2 == 0 {
		return 0
	}
	if l1 < l2 || l1 == 0 {
		return -1
	}

	for i := 0; i <= l1-l2; i++ {
		if needle == haystack[i:i+l2] {
			return i
		}
	}

	return -1
}

func printNumbers(n int) []int {
	l := 0
	res := []int{}
	for 0 < n {
		n--
		l = l*10 + 9
	}
	for i := 0; i < l+1; i++ {
		res = append(res, i)
	}

	return res
}

func isPalindrome(s string) bool {
	var sgood string
	for i := range s {
		if isalnum(s[i]) {
			sgood += string(s[i])
		}
	}

	n := len(sgood)
	sgood = strings.ToLower(sgood)
	for i := 0; i < n/2; i++ {
		if sgood[i] != sgood[i-1-i] {
			return false
		}
	}

	return true
}

func isalnum(char byte) bool {
	return (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') || (char >= 0 && char <= 9)
}

func rotateString(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	var tmp string = a + a
	return strings.Contains(tmp, b)
}

func lengthOfLastWord(s string) int {
	s = strings.Trim(s, " ")
	ss := strings.Split(s, " ")
	if len(ss) == 0 {
		return 0
	}
	return len(ss[len(ss)-1])
}

func lengthOfLastWord2(s string) int {
	tail := len(s) - 1
	for tail >= 0 && s[tail] != ' ' {
		tail--
	}
	if tail < 0 {
		return 0
	}
	head := tail
	for head >= 0 && s[head] != ' ' {
		head--
	}

	return head - tail
}




