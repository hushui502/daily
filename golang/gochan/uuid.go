package gochan

var gochanUUID int

func defaultUUID() int {
	gochanUUID += 1
	return gochanUUID
}
