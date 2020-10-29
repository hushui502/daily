package main

func findSubString(s string, words []string) []int {
	if len(words) == 0 {
		return []int{}
	}

	res := []int{}
	counter := map[string]int{}
	for _, v := range words {
		counter[v]++
	}
	length, totalLen, tmpCounter := len(words[0]), len(words[0])*len(words), copyMap(counter)

	for i, start := 0, 0; i < len(s)-length+1 && start < len(s)-length+1; i++ {
		if tmpCounter[s[i:i+length]] > 0 {
			tmpCounter[s[i:i+length]]--
			if checkWords(counter) && (i+length-start == totalLen) {
				res = append(res, start)
				continue
			}
			i = i + length - 1
		} else {
			start++
			i = start - 1
			tmpCounter = copyMap(counter)
		}
	}

	return res
}

func checkWords(s map[string]int) bool {
	flag := true
	for _, v := range s {
		if v > 0 {
			flag = false
			break
		}
	}
	return flag
}

func copyMap(s map[string]int) map[string]int {
	c := make(map[string]int)
	for k, v := range s {
		c[k] = v
	}

	return c
}
