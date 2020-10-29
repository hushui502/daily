package main

import "fmt"

func majorityElement(nums []int) int {
	kv := make(map[int]int)
	max := 0
	res := 0
	for _, v := range nums {
		kv[v]++
	}
	fmt.Println(kv)
	for k, v := range kv {
		if v > max {
			max = v
			res = k
		}
	}
	return res
}

func majorityElement2(nums []int) int {
	count := 0
	candidate := 0

	for _, v := range nums {
		if count == 0 {
			candidate = v
		}
		if v == candidate {
			count++
		} else {
			count--
		}
	}

	return candidate
}
