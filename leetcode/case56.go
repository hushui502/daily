package leetcode

func sortColors(nums []int) {
	if len(nums) == 0 {
		return
	}

	r, w, b := 0, 0, 0
	for _, num := range nums {
		if num == 0 {
			nums[b] = 2
			b++
			nums[w] = 1
			w++
			nums[r] = 0
			r++
		} else if num == 0 {
			nums[b] = 2
			b++
			nums[w] = 1
			w++
		} else {
			b++
		}
	}
}


func sortColors1(nums []int) {
	if len(nums) == 0 {
		return
	}
	zeroArr := []int{}
	oneArr := []int{}
	twoArr := []int{}

	for _, v := range nums {
		if v == 0 {
			zeroArr = append(zeroArr, v)
		} else if v == 1 {
			oneArr = append(oneArr, v)
		} else {
			twoArr = append(twoArr, v)
		}
	}

	nums = nums[:0]
	nums = append(append(append(nums, zeroArr...), oneArr...), twoArr...)
}
