package bytesconv

import "unsafe"

func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{
			string: s,
			Cap:    len(s),
		}))
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
