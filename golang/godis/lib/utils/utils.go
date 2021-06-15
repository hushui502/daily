package utils

func ToCmdLine(cmd ...string) [][]byte {
	args := make([][]byte, len(cmd))
	for i, s := range cmd {
		args[i] = []byte(s)
	}

	return args
}

func ToCmdLine2(commandName string, args ...string) [][]byte {
	result := make([][]byte, len(args)+1)
	result[0] = []byte(commandName)
	for i, s := range args {
		result[i+1] = []byte(s)
	}
	return result
}

func Equals(a, b interface{}) bool {
	sliceA, okA := a.([]byte)
	sliceB, oKB := b.([]byte)
	if okA && oKB {
		return BytesEquals(sliceA, sliceB)
	}
	return a == b
}

func BytesEquals(a, b []byte) bool {
	if a == nil && b == nil {
		return true
	}
	if a != nil && b == nil {
		return false
	}
	if a == nil && b != nil {
		return false
	}

	size := len(a)
	for i := 0; i < size; i++ {
		av := a[i]
		bv := b[i]
		if av != bv {
			return false
		}
	}

	return true
}
