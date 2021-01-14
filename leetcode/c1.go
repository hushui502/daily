package leetcode

func interset(nums1, nums2 []int) []int {
	m0 := make(map[int]int)
	for _, v := range nums1 {
		m0[v] += 1
	}
	k := 0
	for _, v := range nums2 {
		if m0[v] > 0 {
			m0[v]--
			nums2[k] = v
			k++
		}
	}

	return nums2[0:k]
}

// sort array
func interset2(nums1, nums2 []int) []int {
	l1 := len(nums1)
	l2 := len(nums2)

	i, j, k := 0, 0, 0

	for i < l1 && j < l2 {
		if nums1[i] < nums2[j] {
			i++
		}
		if nums1[i] > nums2[j] {
			j++
		}
		if nums1[i] == nums2[j] {
			// eg. nums1[0] = nums1[i] 原地修改,节约内存
			nums1[k] = nums1[i]
			i++
			j++
			k++
		}
	}

	return nums1[:n]
}
