package main

import "fmt"

//func deleteElement(s []int, element int) []int {
//	for i := 0; i < len(s); i++ {
//		if s[i] == element {
//			s = append(s[:i], s[i+1:]...)
//			println(s)
//			fmt.Println(s)
//			//s = s[:i]
//			//s = append(s, s[i+1:]...)
//			return s
//		}
//	}
//	return s
//
//}
//
//func main() {
//	a := []int{1, 2, 3, 4}
//	println(a)
//	a = deleteElement(a, 2)
//	println(a)
//	fmt.Println(a)
//}

func main() {
	s := []int{1, 2, 3, 4}
	for i := 0; i < len(s); i++ {
		if s[i] == 2 {
			s = append(s[:i], s[i+1:]...)
			break
		}
	}
	fmt.Println(s)
}
