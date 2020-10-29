package pattern

import "fmt"

type Handler interface {
	handle(content string)
	next(handler Handler, content string)
}

type A struct {
	handler Handler
}

func (a *A) handle(content string) {
	fmt.Println("a handle")
	nextContent := ""
	a.next(a.handler, nextContent)
	
}

func (a *A) next(handler Handler, content string) {
	if a.handler != nil {
		a.handler.handle(content)
	}
}

type B struct {
	handler Handler
}

func (b *B) handle(content string) {
	fmt.Println("a handle")
	nextContent := ""
	b.next(b.handler, nextContent)

}

func (b *B) next(handler Handler, content string) {
	if b.handler != nil {
		b.handler.handle(content)
	}
}

func testchain() {
	aHandler := &A{}
	bHandler := &B{}

	aHandler.handler = bHandler

	aHandler.handle("txt")
}


