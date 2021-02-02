package main

//func main() {
//	// 0
//	// 空结构体不占内存空间，仅仅作为占位符十分适合。
//	fmt.Println(unsafe.Sizeof(struct{}{}))
//}

// Set实现
// Go中没有Set的实现，所以手动实现的话考虑map替代，但是只需要用到map的键，不需要值。
// 所以使用struct{}要比bool好很多，因为使用bool每个需要1个字节的空间占用。

type Set map[string]struct{}

func (s Set) Has(key string) bool {
	_, ok := s[key]
	return ok
}

func (s Set) Add(key string) {
	s[key] = struct{}{}
}

func (s Set) Delete(key string) {
	delete(s, key)
}

// 不发送数据的信道
// 很多channel只是起到消息通知的目的，但是需要send，recv一些信息才可以

func write(ch chan struct{}) {
	<-ch
	println("```")
	close(ch)
}

func main() {
	ch := make(chan struct{})
	go write(ch)
	ch <- struct{}{}
}

// 在部分场景下，结构体只包含方法，不包含任何的字段。
// 这里也可以使用int，float等替代，但是会额外的内存消耗
type Student struct {}

func (s Student) Go() {

}

func (s Student) Back() {

}