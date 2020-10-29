package main

import (
	"strconv"
	"strings"
)

func compressString(s string) string {
	if len(s) == 0 {
		return s
	}
	ans := ""
	ch := s[0]
	cnt := 1
	for i := 1; i < len(s); i++ {
		if ch == s[i] {
			cnt++
		} else {
			ans += string(ch) + strconv.Itoa(cnt)
			cnt = 1
			ch = s[i]
		}
	}
	ans += string(ch) + strconv.Itoa(cnt)
	if len(ans) < len(s) {
		return s
	}
	return ans
}

func method2(s string) string {
	if s == "" {
		return ""
	}
	var stringBuilder strings.Builder
	curr := s[0]
	curLen := 1

	for i := 0; i < len(s); i++ {
		if s[i] == curr {
			curLen++
		} else {
			stringBuilder.WriteByte(curr)
			stringBuilder.WriteString(strconv.Itoa(curLen))
			curr = s[i]
			curLen = 1
		}
	}
	stringBuilder.WriteByte(curr)
	stringBuilder.WriteString(strconv.Itoa(curLen))

	if stringBuilder.Len() >= len(s) {
		return s
	}
	return stringBuilder.String()
}

func compressString1(S string) string {
	if len(S) <= 1 {
		return S
	}
	res := ""
	slow, quick := 0, 0
	for quick < len(S) {
		if S[quick] != S[slow] {
			res += string(S[slow]) + strconv.Itoa(quick-slow)
			slow = quick
		}
		quick++
	}

	res += string(S[slow]) + strconv.Itoa(quick-slow)

	if len(res) >= len(S) {
		return S
	}
	return res
}

type stu struct {
	name string
}

func test() {
	m := map[int]stu{
		1: {name: "hufan"},
	}
	//a1 := m[1].name

}
