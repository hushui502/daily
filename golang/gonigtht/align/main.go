package align

import (
	"fmt"
	"sync/atomic"
)

func main() {
	var c0 int64
	fmt.Println(atomic.AddInt64(&c0, 1))

	c1 := [5]int64{}
	fmt.Println(atomic.AddInt64(&c1[:][0], 1))
}
