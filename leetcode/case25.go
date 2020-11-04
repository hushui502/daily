package main

func search33(nums []int, target int) int {
	if len(nums) == 0 {
		return -1
	}

	low, high := 0, len(nums)-1
	for low < high {
		mid := low + (high-low)>>1
		if nums[mid] == target {
			return mid
		} else if nums[low] < nums[mid] {
			if nums[low] <= target && target < nums[mid] {
				high = mid - 1
			} else {
				low = mid + 1
			}
		} else if nums[mid] < nums[high] {
			if nums[mid] < target && target <= nums[high] {
				low = mid + 1
			} else {
				high = mid - 1
			}
		} else {
			if nums[low] == nums[mid] {
				low++
			}
			if nums[mid] == nums[high] {
				high--
			}
		}
	}

	return -1
}
