package main

import "testing"

func TestSafeEqual(t *testing.T) {
	str1 := "abc"
	str2 := "abc"
	str3 := "ab3"

	if !safeEqual(str1, str2) {
		t.Fatal("error ...")
	}
	if safeEqual(str1, str3) {
		t.Fatal("error ..")
	}
}
