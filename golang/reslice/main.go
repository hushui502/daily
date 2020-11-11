package main

import "fmt"

// ref: https://mp.weixin.qq.com/s/VVB1-7DcaYmbcppK7gDx_A


func c() []int {
	a1 := []int{1}
	// no panic LOL
	a2 := a1[1:]

	return a2
}

func main() {
	res := c()
	fmt.Println(res)

	outerSlice := []string{"a", "a"}
	fmt.Printf("%p %v %p\n", &outerSlice, outerSlice, &outerSlice)
	modifySlice(outerSlice)
	fmt.Printf("%p %v %p\n", &outerSlice, outerSlice, &outerSlice)
}


func modifySlice(innerSlice []string) {
	fmt.Printf("%p %v %p\n", &innerSlice, innerSlice, &innerSlice[0])
	// 扩容会导底层数组引用改变，所以要看是否会扩容
	innerSlice = append(innerSlice, "a")
	innerSlice[0] = "b"
	innerSlice[1] = "b"
	fmt.Printf("%p %v %p\n", &innerSlice, innerSlice, &innerSlice[0])
}
