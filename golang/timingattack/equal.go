package main

// Constant time for same length String comparison, to prevent timing attacks
func safeEqual(a, b string) bool {
	if len(a) != len(b) {
		return false
	}

	var equal uint8
	for i := 0; i < len(a); i++ {
		equal |= a[i]^b[i]
	}

	return equal == 0
}

func unsafeEqual(a, b string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}


func main() {
	a := "abc"
	b := "abc"
	println(unsafeEqual(a, b))
}
