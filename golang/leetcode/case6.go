package main

import "strconv"

func isPalindrome(x int) bool {
	if x < 0 {
		return false
	}
	if x < 10 {
		return false
	}

	s := strconv.Itoa(x)
	len := len(s)
	for i := 0; i < len/2; i++ {
		if s[i] != s[len-1-i] {
			return false
		}
	}

	return true
}
