package leetcode

func plusOne(digits []int) []int {
	var res []int
	addOn := 0
	for i := len(digits)-1; i >= 0; i-- {
		digits[i] += addOn
		// 这里是为了避免对addOn多次增加,其实这题的简单之处在于只加1,如果是别的数我们就需要额外的判断
		addOn = 0
		if i == len(digits)-1 {
			digits[i] += 1
		}
		if digits[i] == 10 {
			addOn = 1
			digits[i] = digits[i] % 10
		}
	}
	if addOn == 1 {
		res = append([]int{1}, digits...)
	} else {
		res = digits
	}

	return res
}
