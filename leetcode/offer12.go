package leetcode

func largeGroupPositions(s string) [][]int {
	cnt := 1
	var ans [][]int
	for i := range s {
		if i == len(s)-1 || s[i] != s[i+1] {
			if cnt >= 3 {
				ans = append(ans, []int{i-cnt+1, i})
			}
			cnt = 1
		} else {
			cnt++
		}
	}

	return ans
}
