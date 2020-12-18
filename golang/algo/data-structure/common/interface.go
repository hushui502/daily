package common

// Compare is a generic interface that represents items that can
// be compared
type Comparator interface {
	// Compare compare this interface with another.
	// Returns a positive number if this interface is greeter
	// 0 if equal, negative number if less
	Compare(Comparator) int
}

// Comparators is a typed list of type Comparator
type Comparators []Comparator
