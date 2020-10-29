package main

func canThreePartsEqualSum(A []int) bool {
	sum := 0
	for _, v := range A {
		sum += v
	}
	if sum%3 != 0 {
		return false
	}
	sum1, sum2 := sum/3, sum/3
	for i := range A {
		if sum1 != 0 {
			sum1 -= A[i]
		} else if sum2 != 0 {
			sum2 -= A[i]
			if sum2 == 0 {
				return true
			}
		}
	}
	return false
}

func Method2(A []int) bool {
	sum := 0
	for i := 0; i < len(A); i++ {
		sum += A[i]
	}

	if sum%3 != 0 {
		return false
	}
	left := 0
	right := len(A) - 1
	leftSum := A[left]
	rightSum := A[right]
	segVal := sum / 3

	for left+1 < right {
		if leftSum == rightSum {
			if leftSum == segVal {
				return true
			}
		}
		if leftSum < sum/3 {
			left++
			leftSum += A[left]
		}
		if rightSum < sum/3 {
			right--
			rightSum += A[right]
		}
	}
	return false
}
