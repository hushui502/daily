package pointer

import (
	"log"
	"testing"
	"unsafe"
)

type puintptr uintptr

//go:nosplit
func (pp puintptr) ptr() *p { return (*p)(unsafe.Pointer(pp)) }

//go:nosplit
func (pp *puintptr) set(p *p) { *pp = puintptr(unsafe.Pointer(p)) }

type p struct {
	id   int
	link puintptr
}

var pidleList puintptr

func pidleput(_p_ *p) {
	_p_.link = pidleList
	pidleList.set(_p_)
}

func pidleget() *p {
	_p_ := pidleList.ptr()
	if _p_ != nil {
		pidleList = _p_.link
	}
	return _p_
}

func TestPtoPList(t *testing.T) {
	for i := 0; i < 5; i++ {
		tp := p{id: i}
		pidleput(&tp)
		log.Printf("put p: %d", i)
	}
	log.Println("")
	for i := 0; i < 5; i++ {
		tp := pidleget()
		log.Printf("get p: %d", tp.id)
	}
}

// TestUintptr1: pointer_test.go:54: a.addr=0xc00000a368, b.addr=0xc00000a368, *b=100
// TestUintptr1: pointer_test.go:56: buintptr=0xc00000a368
func TestUintptr1(t *testing.T) {
	var a int
	a = 100
	b := (*int)(unsafe.Pointer(&a))
	t.Logf("a.addr=%p, b.addr=%p, *b=%d", &a, b, *b)
	buintptr := uintptr(unsafe.Pointer(&a))
	t.Logf("buintptr=0x%x", buintptr)
}

// TestUintptr2: pointer_test.go:65: a0.addr=0xc00000a370, a1.addr=0xc00000a378
// TestUintptr2: pointer_test.go:67: auintptr=0xc00000a370
// TestUintptr2: pointer_test.go:71: a1pointer=0xc00000a378, a1=2
func TestUintptr2(t *testing.T) {
	var a [2]int
	a[0] = 1
	a[1] = 2

	t.Logf("a0.addr=%p, a1.addr=%p", &a[0], &a[1])
	auintptr := uintptr(unsafe.Pointer(&a))
	t.Logf("auintptr=0x%x", auintptr)
	a1uintptr := auintptr + unsafe.Sizeof(&a)
	a1pointer := unsafe.Pointer(a1uintptr)
	a1 := *(*int)(a1pointer)
	t.Logf("a1pointer=%p, a1=%d", a1pointer, a1)

}
