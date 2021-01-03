package leetcode

import "strings"

func replaceSpace(s string) string {
	// maybe realloc
	var b []byte
	for _, v := range s {
		if v == 32 {
			b =append(b,37,50,48)
		} else {
			b = append(b, v)
		}
	}

	return string(b)
}

func replaceSpace2(s string) string {
	var res strings.Builder
	for i := range s {
		if s[i] == ' ' {
			res.WriteString("%20")
		} else {
			res.WriteByte(s[i])
		}
	}

	return res.String()
}

func replaceSpace3(s string) string {
	return strings.ReplaceAll(s, " ", "%20")
}





