package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

/*
	slice reallocated
*/
//func main() {
//	path := []byte("AAAA/BBBBB")
//	sepIndex := bytes.IndexByte(path, '/')
//
//	dir1 := path[:sepIndex:sepIndex]
//	dir2 := path[sepIndex+1:]
//
//	fmt.Println("dir1 => ", string(dir1))
//	fmt.Println("dir2 => ", string(dir2))
//
//
//	dir1 = append(dir1, "suffix"...)
//	fmt.Println("dir1 => ", string(dir1))
//	fmt.Println("dir2 => ", string(dir2))
//}


/*
	deep comparison
*/
//type data struct {
//	num int
//	check [10]func() bool
//	doit func() bool
//	m map[string]string
//	bytes []byte
//}
//
//func main() {
//	v1 := data{}
//	v2 := data{}
//	fmt.Println("deep equal v1 == v2 ? ", reflect.DeepEqual(v1, v2))
//
//	m1 := map[string]string{"one":"a", "two":"b"}
//	m2 := map[string]string{"two":"b", "one":"a"}
//	fmt.Println("v1 == v2 ? ", reflect.DeepEqual(m1, m2))
//
//	s1 := []int{1, 2, 3}
//	s2 := []int{1, 2, 3}
//	fmt.Println("s1 == s2 ? ", reflect.DeepEqual(s1, s2))
//}


/*
	Interface pattern
*/
//type Country struct {
//	WithName
//}
//
//type City struct {
//	WithName
//}
//
//type WithName struct {
//	Name string
//}
//
//type Printable interface {
//	PrintStr()
//}
//
//func (w WithName) PrintStr() {
//	fmt.Println(w.Name)
//}
//
//func main() {
//	city := City{WithName{"BEIJING"}}
//	country := Country{WithName{"CHINA"}}
//	city.PrintStr()
//	country.PrintStr()
//}

/*
	var AInterface = (*AImpl)(nil)
 */
//type Shape interface {
//	Sides() int
//	Area() int
//}
//
//type Square struct {
//	len int
//}
//
//func (s *Square) Area() int {
//	return s.len * s.len
//}
//
//func (s *Square) Sides() int {
//	return 4
//}
//
//func main() {
//	var _ Shape = (*Square)(nil)
//	s := Square{len:4}
//	fmt.Printf("%d\n", s.Sides())
//}


/*
	Delegation e.g1
*/
//type Widget struct {
//	X, Y int
//}
//
//type Label struct {
//	Widget
//	Text string
//}
//
//func (label Label) Paint() {
//	fmt.Println("%p: Label.Paint(%q)\n", &label, label.Text)
//}
//
//type Painter interface {
//	Paint()
//}
//
//type Clicker interface {
//	Click()
//}
//
//type Button struct {
//	Label
//}
//
//func NewButton(x, y int, text string) Button {
//	return Button{Label{Widget{x, y}, text}}
//}
//
//func (button Button) Paint() {
//	fmt.Println("Button.Paint(%s)\n", button.Text)
//}
//
//func (button Button) Click() {
//	fmt.Println("Button.Click(%s)\n", button.Text)
//}
//
//type ListBox struct {
//	Widget
//	Texts []string
//	Index int
//}
//
//func (listBox ListBox) Paint() {
//	fmt.Println("ListBox.Paint(%q)\n", listBox.Texts)
//}
//
//func (listBox ListBox) Click() {
//	fmt.Println("ListBox.Click(%q)\n", listBox.Texts)
//}
//
//func main() {
//	button1 := Button{Label{Widget{10, 10}, "OK"}}
//	button2 := Button{Label{Widget{20, 20}, "Cancel"}}
//	listBox := ListBox{Widget{13, 23},
//		[]string{"AL", "AZ", "AR", "AK"}, 0}
//
//	for _, painter := range []Painter{listBox, button1, button2} {
//		painter.Paint()
//	}
//
//	for _, widget := range []interface{}{listBox, button1, button2} {
//		if clicker, ok := widget.(Clicker); ok {
//			clicker.Click()
//		}
//	}
//}

/**
	err handle 1
 */
func parse(r io.Reader) {
	var err error
	read := func(data interface{}) {
		if err != nil {
			return
		}
		_, err = ioutil.ReadFile(data.(string))
	}

	read("/xx/sss")
	if err != nil {
		return
	}
}


/**
	err handle 2
 */
type Reader struct {
	r io.Reader
	err error
}

func (r *Reader) read(data interface{}) {
	if r.err != nil {
		_, r.err = ioutil.ReadFile(data.(string))
	}
}

func parse1(input io.Reader) {
	r := Reader{r:input}
	r.read("hello")

	if r.err != nil {
		return
	}
}


/**
	functional option
 */
//type Server struct {
//	Addr string
//	Port int
//	Protocol string
//	Timeout time.Duration
//	MaxCounts int
//	TLS *tls.Config
//}
//
//type Option func(*Server)
//
//func Protocol(p string) Option {
//	return func(server *Server) {
//		server.Protocol = p
//	}
//}
//
//func Timeout(timeout time.Duration) Option {
//	return func(server *Server) {
//		server.Timeout = timeout
//	}
//}
//
//func MaxCounts(maxCounts int) Option {
//	return func(server *Server) {
//		server.MaxCounts = maxCounts
//	}
//}
//
//func TLS(tls *tls.Config) Option {
//	return func(server *Server) {
//		server.TLS = tls
//	}
//}
//
//func NewServer(addr string, port int, options ...func(*Server)) (*Server, error) {
//	// you can setup a default config
//	srv := Server{
//		Addr:addr,
//		Port:port,
//	}
//	for _, option := range options {
//		option(&srv)
//	}
//
//	return &srv, nil
//}
//
//func main() {
//	s1, _ := NewServer("localhost", 1111)
//	s2, _ := NewServer("localhost", 2222, Protocol("tcp"))
//	fmt.Println(s1)
//	fmt.Println(s2)
//}


/**
	map/reduce/filter
 */
// map reduce 1
func MapUpCase(arr []string, fn func(s string) string) []string {
	var newArray = []string{}
	for _, item := range arr {
		newArray = append(newArray, fn(item))
	}

	return newArray
}

func MapLen(arr []string, fn func(s string) int) []int {
	var newArray = []int{}
	for _, item := range arr {
		newArray = append(newArray, fn(item))
	}

	return newArray
}

func Reduce(arr []string, fn func(s string) int) int {
	sum := 0
	for _, it := range arr {
		sum += fn(it)
	}
	
	return sum
}

func Filter(arr []int, fn func(n int) bool) []int {
	var newArray = []int{}
	for _, it := range arr {
		if fn(it) {
			newArray = append(newArray, it)
		}
	}

	return newArray
}

type Employee struct {
	Name string
	Age int
	Vacation int
	Salary int
}

var employeeList = []Employee{
	{"hufan", 22, 10, 22222},
	{"bob", 32, 23, 100203},
	{"alice", 44, 33, 322222},
}

func EmployeeCountIf(list []Employee, fn func(e *Employee) bool) int {
	count := 0
	for i, _ := range list {
		if fn(&list[i]) {
			count += 1
		}
	}

	return count
}

/**
	generic map
 */
func transform(slice, function interface{}, inPalce bool) interface{} {
	sliceInType := reflect.ValueOf(slice)
	if sliceInType.Kind() != reflect.Slice {
		panic("not slice")
	}

	fn := reflect.ValueOf(function)
	elemType := sliceInType.Type().Elem()
	if !vertifyFuncSignature(fn, elemType, nil) {
		panic("transform: function must be of type func(" + sliceInType.Type().Elem().String() + ") outputElementType")
	}

	sliceOutType := sliceInType
	if !inPalce {
		sliceOutType = reflect.MakeSlice(reflect.SliceOf(fn.Type().Out(0)), sliceInType.Len(), sliceInType.Len())
	}
	for i := 0; i < sliceInType.Len(); i++ {
		sliceOutType.Index(i).Set(fn.Call([]reflect.Value{sliceInType.Index(i)})[0])
	}

	return sliceOutType.Interface()
}

func vertifyFuncSignature(fn reflect.Value, types ...reflect.Type) bool {
	if fn.Kind() != reflect.Func {
		return false
	}

	if (fn.Type().NumIn() != len(types)-1) || (fn.Type().NumOut() != 1) {
		return false
	}

	for i := 0; i < len(types)-1; i++ {
		if fn.Type().In(i) != types[i] {
			return false
		}
	}

	outType := types[len(types)-1]
	if outType != nil && fn.Type().Out(0) != outType {
		return false
	}

	return true
}

func Transform(slice, fn interface{}) interface{} {
	return transform(slice, fn ,false)
}

func TransformInPlace(slice, fn interface{}) interface{} {
	return transform(slice, fn, true)
}

/*
	generic reduce
*/
func Reduce1(slice, pairFunc, zero interface{}) interface{} {
	sliceInType := reflect.ValueOf(slice)
	if sliceInType.Kind() != reflect.Slice {
		panic("reduce: wrong type, not slice")
	}

	len := sliceInType.Len()
	if len == 0 {
		return zero
	} else if len == 1 {
		return sliceInType.Index(0)
	}

	elemType := sliceInType.Type().Elem()
	fn := reflect.ValueOf(pairFunc)
	if !vertifyFuncSignature(fn, elemType, elemType, elemType) {
		t := elemType.String()
		panic("reduce: function must be of type func(" + t + ", " + t + ") " + t)
	}

	var ins [2]reflect.Value
	ins[0] = sliceInType.Index(0)
	ins[1] = sliceInType.Index(1)
	out := fn.Call(ins[:])[0]

	for i := 2; i < len; i++ {
		ins[0] = out
		ins[1] = sliceInType.Index(i)
		out = fn.Call(ins[:])[0]
	}

	return out.Interface()
}

func mul(a, b int) int {
	return a * b
}

/**
	generic filter
 */
var boolType = reflect.ValueOf(true).Type()

func filter1(slice, function interface{}, inPlace bool) (interface{}, int) {
	sliceInType := reflect.ValueOf(slice)
	if sliceInType.Kind() != reflect.Slice {
		panic("filter: wrong type, not a slice")
	}

	fn := reflect.ValueOf(function)
	elemType := sliceInType.Type().Elem()
	if !vertifyFuncSignature(fn, elemType, boolType) {
		panic("filter: function must be of type func(" + elemType.String() + ") bool")
	}

	var which []int
	for i := 0; i < sliceInType.Len(); i++ {
		if fn.Call([]reflect.Value{sliceInType.Index(i)})[0].Bool() {
			which = append(which, i)
		}
	}

	out := sliceInType

	if !inPlace {
		out = reflect.MakeSlice(sliceInType.Type(), len(which), len(which))
	}
	for i := range which {
		out.Index(i).Set(sliceInType.Index(which[i]))
	}

	return out.Interface(), len(which)
}

func Filter1(slice, fn interface{}) interface{}  {
	result, _ := filter1(slice, fn, false)
	return result
}

func FilterInPlace1(slicePtr, fn interface{}) {
	in := reflect.ValueOf(slicePtr)
	if in.Kind() != reflect.Ptr {
		panic("FilterInPlace: wrong type, " + "not a pointer to slice")
	}
	_, n := filter1(in.Elem().Interface(), fn, false)
	in.Elem().SetLen(n)
}

func isEven(a int) bool {
	return a%2 == 0
}

func isOddString(s string) bool {
	i, _ := strconv.ParseInt(s, 10, 32)
	return i%2 == 1
}


func main() {
	var list = []string{"heLlO", "Chen", "Hu"}

	// map
	x := MapUpCase(list, func(s string) string {
		return strings.ToUpper(s)
	})
	fmt.Printf("%v\n", x)

	y := MapLen(list, func(s string) int {
		return len(s)
	})
	fmt.Printf("%v\n", y)

	// reduce
	xlen := Reduce(list, func(s string) int {
		return len(s)
	})
	fmt.Printf("%v\n", xlen)

	// filter
	var inset = []int{1, 2, 3, 4, 5, 6}
	out := Filter(inset, func(n int) bool {
		return n % 2 == 1
	})
	fmt.Printf("%v\n", out)

	// generic map
	old := EmployeeCountIf(employeeList, func(e *Employee) bool {
		return e.Age > 30
	})
	fmt.Printf("old people(>30): %v\n", old)

	// generic reduce
	a := make([]int, 10)
	for i := range a {
		a[i] = i + 1
	}
	genericReduceOut := Reduce1(a, mul, 1).(int)
	fmt.Printf("genericReduceOut: %v\n", genericReduceOut)

	// generic filter
	a1 := []int{1, 2, 3, 4}
	result := Filter1(a1, isEven)
	fmt.Printf("%v\n", result)

	s1 := []string{"1", "2", "3", "4"}
	result = Filter1(s1, isOddString)
	fmt.Printf("%v\n", result)
}



















