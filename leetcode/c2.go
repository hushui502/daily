package leetcode

import "strings"

func longestCommandPrefix(strs []string) string {
	if len(strs) < 1 {
		return ""
	}

	prefix := strs[0]
	for _, v := range strs {
		// 这里利用index不存在的话则返回-1
		// 如果是返回0则说明是当前的prefix
		for strings.Index(v, prefix) != 0 {
			if len(prefix) == 0 {
				return ""
			}
			prefix = prefix[:len(prefix)-1]
		}
	}

	return prefix
}
