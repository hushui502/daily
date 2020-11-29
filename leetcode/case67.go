package leetcode

func merge(nums1 []int, m int, nums2 []int, n int) {
	if m == 0 {
		copy(nums1, nums2)
		return
	}

	i := m - 1
	j := n - 1
	k := i + j -1

	for ; i >= 0 && j >= 0; k-- {
		if nums1[i] > nums2[j] {
			nums1[k] = nums1[i]
			i--
		} else {
			nums1[k] = nums2[j]
			j--
		}
	}

	for ; j >= 0; k-- {
		nums1[k] = nums2[j]
		j--
	}
}
