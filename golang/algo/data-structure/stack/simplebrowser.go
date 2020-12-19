package stack


// simple browser
// that has many problems (LOL
type Browser struct {
	forwardStack Stack
	backStack Stack
}

func NewBrowser() *Browser {
	return &Browser{
		forwardStack: NewArrayStack(),
		backStack:    NewLinedListStack(),
	}
}

func (bs *Browser) CanFroward() bool {
	if bs.forwardStack.IsEmpty() {
		return false
	}
	return true
}

func (bs *Browser) CanBack() bool {
	if bs.backStack.IsEmpty() {
		return false
	}

	return true
}

func (bs *Browser) Open(addr string) {
	// logic...
	bs.forwardStack.Flush()
}

func (bs *Browser) PushBack(addr string) {
	bs.backStack.Push(addr)
}

func (bs *Browser) Forward() {
	if bs.forwardStack.IsEmpty() {
		return
	}
	top := bs.forwardStack.Pop()
	bs.backStack.Push(top)
}

func (bs *Browser) Back() {
	if bs.backStack.IsEmpty() {
		return
	}

	top := bs.backStack.Pop()
	bs.forwardStack.Push(top)
}
