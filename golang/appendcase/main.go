package main

import (
	"fmt"
)
// slice 扩容会取消对原数组的底层引用，因此，一般情况我们用copy比较好，且最好指定slice的可控容量
func main() {
	ar := make([]int64, 2, 2)
	fmt.Println(ar)

	arDeriv1 := ar
	arDeriv1 = append(arDeriv1, 0)
	fmt.Println(ar, arDeriv1)

	arDeriv2 := ar
	arDeriv2 = append(arDeriv2, 1)
	fmt.Println(ar, arDeriv1, arDeriv2)
}