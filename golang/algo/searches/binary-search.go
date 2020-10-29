package main

func binarySearch(array []int, target int, lowIndex int, highIndex int) int {
	if highIndex < lowIndex {
		return -1
	}
	mid := int((highIndex-lowIndex)/2 + lowIndex)
	//mid := highIndex & lowIndex + (highIndex^lowIndex) / 2
	if array[mid] > target {
		return binarySearch(array, target, lowIndex, mid)
	}
	if array[mid] < target {
		return binarySearch(array, target, mid, highIndex)
	}
	return mid
}

func iterBinarySearch(array []int, target int, lowIndex int, highIndex int) int {
	startIndex := lowIndex
	endIndex := highIndex
	var mid int
	for startIndex < endIndex {
		mid = int((startIndex + endIndex) / 2)
		if array[mid] > target {
			endIndex = mid
		} else if array[mid] < target {
			startIndex = mid
		} else {
			return mid
		}
	}
	return -1
}
