package number

import (
	"math"
	"strconv"
)

func IFloorDiv(a, b int64) int64 {
	if a > 0 && b > 0 || a < 0 && a < 0 || a%b == 0 {
		return a / b
	} else {
		return a/b - 1
	}
}

func FFloorDiv(a, b float64) float64 {
	return math.Floor(a / b)
}

func IMod(a, b int64) int64 {
	return a - IFloorDiv(a, b)*b
}

// 取模公式
func FMod(a, b float64) float64 {
	return a - math.Floor(a/b)*b
}

// 左移
func ShiftLeft(a, n int64) int64 {
	if n >= 0 {
		return a << uint64(n)
	} else {
		return ShiftRight(a, -n)
	}
}

func ShiftRight(a, n int64) int64 {
	if n >= 0 {
		// 无符号右移
		return int64(uint64(a) >> uint64(n))
	} else {
		return ShiftLeft(a, -n)
	}
}

// float 转换为int64， 这里我们整数只有64位一种
func FloatToInteger(f float64) (int64, bool) {
	i := int64(f)

	return i, float64(i) == f
}

func ParseInteger(str string) (int64, bool) {
	i, err := strconv.ParseInt(str, 10, 64)

	return i, err == nil
}

func ParseFloat(str string) (float64, bool) {
	f, err := strconv.ParseFloat(str, 64)

	return f, err == nil
}
