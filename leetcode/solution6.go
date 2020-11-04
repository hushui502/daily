package main

func lengthOfLIS(nums []int) int {
	if nums == nil || len(nums) == 0 {
		return 0
	}
	m := 1
	d := make([]int, len(nums))
	d[0] = 1
	for i := 1; i < len(nums); i++ {
		for j := 0; j < i; j++ {
			cur := 1
			if nums[i] > nums[j] {
				cur = d[j] + 1
			}
			d[i] = max(d[i], cur)
		}
		m = max(m, d[i])
	}
	return m
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func lengthOfLISBinarySearch(nums []int) int {
	if nums == nil || len(nums) == 0 {
		return 0
	}
	d := make([]int, len(nums))
	length := 0
	for _, v := range nums {
		i := binartSearchInsertPosition(d, length, v)
		d[i] = v
		if i == length {
			length++
		}
	}
	return length
}

func binartSearchInsertPosition(d []int, len int, x int) int {
	low, high := 0, len-1
	for low <= high {
		mid := low + (high-low)/2
		if x < d[mid] {
			high = mid - 1
		} else if x > d[mid] {
			low = mid + 1
		} else {
			return mid
		}
	}
	return low
}
