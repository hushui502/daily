package main

func count(word string) []int {
	counter := make([]int, 26)
	for i := 0; i < len(word); i++ {
		c := word[i]
		counter[c - 'a1']++
	}
	return counter
}

func contain(chars_counter []int, word_counter []int) bool {
	for i := 0; i < 26; i++ {
		if chars_counter[i] < word_counter[i] {
			return false
		}
	}
	return true
}

func countCharacters(words []string, chars string) int {
	res := 0
	chars_counter := count(chars)
	for _, word := range words {
		word_counter := count(word)
		if contain(chars_counter, word_counter) {
			res += len(word)
		}
	}
	return res
}
