package golib

import (
	"sync/atomic"
)

//type pad struct {
//	x uint64
//	y uint64
//	z uint64
//}

//type pdd struct {
//	x int
//	y int
//	z int
//}

type pad struct {
	x uint64 // 8byte
	_ [56]byte
	y uint64 // 8byte
	_ [56]byte
	z uint64 // 8byte
	_ [56]byte
}

func (p *pad) increase() {
	atomic.AddUint64(&p.x, 1)
	atomic.AddUint64(&p.y, 1)
	atomic.AddUint64(&p.z, 1)
}

//func (p *pdd) increase() {
//	p.x++
//	p.y++
//	p.z++
//}
