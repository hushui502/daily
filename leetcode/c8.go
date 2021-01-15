package leetcode

import "strings"

func convert(s string, numRows int) string {
	if numRows == 1 {
		return s
	}

	b := []rune(s)
	res := make([]string, numRows)
	length := len(b)
	period := numRows * 2 - 2
	for i := 0; i < length; i++ {
		mod := i % period
		if mod < numRows {
			res[i] += string(b[i])
		} else {
			res[period - mod] += string(b[i])
		}
	}

	return strings.Join(res, "")
}
