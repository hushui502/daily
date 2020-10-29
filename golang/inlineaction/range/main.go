package main

const size = 1000000

type SomeStruct struct {
	ID0 int64
	ID1 int64
	ID2 int64
	ID3 int64
	ID4 int64
	ID5 int64
	ID6 int64
	ID7 int64
	ID8 int64
	ID9 int64
}

func main() {
	const size = 1000000

	slice := make([]SomeStruct, size)
	for i := 0; i < len(slice); i++ {
		_ = slice[i]
	}
	//for _, s := range slice {
	//	_ = s
	//}
}
