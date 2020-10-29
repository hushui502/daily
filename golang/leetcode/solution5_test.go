package main

import "testing"

func TestMethod2(t *testing.T) {
	arr := []int{2, 2, 1, 1, 1, 2, 2}
	a := majorityElement(arr)

	if a != 2 {
		t.Fatal("failed")
	}
}
