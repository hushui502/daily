package reflect

import (
	"fmt"
	"testing"
)

func TestMakeSum(t *testing.T) {
	var intSum func(int, int) int64

	MakeSum(intSum)

	fmt.Println(intSum(1, 2))
}
