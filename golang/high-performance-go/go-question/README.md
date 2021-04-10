## 前言
虽然都是些小问题，但都可以延申，其实有了答案的问题本身就算不上高深的问题。至于有了答案仍然看不懂学不会的问题，那肯定是答案的问题和你理解的问题，答案的问题主要存在于解答者自身的理解和叙述存在问题，理解的问题主要在于求解者思路的问题，其实两者的核心都在于自身的知识储备问题，遇到问题，可以尝试去从一些基础的概念入手，最后组织推导问题的答案，大多是可以解决的。遇到越来越多的问题，是件好事，证明还在进步，对自己有要求，
## 问题
-  = 和 := 的区别？
```
:= 在golang中是声明+赋值
= 只是单纯的赋值
var a int
a = 10

b := 10

声明的作用是什么呢？如果没有声明，没有类型这个概念，会有什么坏处和好处。本质上不都是内存上的地址，我们的声明是在语法词法分析还是在地址中有空间指明？
```
- 指针的作用？
```
指针像是一种类型，指针并不是地址。
我们主要理解清楚 类型，变量，值三者之间的关系就可以了。
var a = 10
var p *int = &a
fmt.Printf("a = %d",  *p) // a 可以用 *p 访问
```
- Go 允许多个返回值吗？
```
允许。
Lua中的多返回值其实完全是根据栈来操作的，

"".test t=1 size=32 value=0 args=0x20 locals=0x0
        0x0000 00000 (test.go:5)        TEXT    "".test(SB), $0-32//栈大小为32字节
        0x0000 00000 (test.go:5)        NOP
        0x0000 00000 (test.go:5)        NOP
        0x0000 00000 (test.go:5)        MOVQ    "".i+8(FP), CX//取第一个参数i
        0x0005 00005 (test.go:5)        MOVQ    "".j+16(FP), AX//取第二个参数j
        0x000a 00010 (test.go:5)        FUNCDATA        $0, gclocals·a8eabfc4a4514ed6b3b0c61e9680e440(SB)
        0x000a 00010 (test.go:5)        FUNCDATA        $1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x000a 00010 (test.go:6)        MOVQ    CX, BX//将i放入bx
        0x000d 00013 (test.go:6)        ADDQ    AX, CX//i+j放入cx
        0x0010 00016 (test.go:7)        SUBQ    AX, BX//i-j放入bx
						//将返回结果存入调用函数栈帧
        0x0013 00019 (test.go:8)        MOVQ    CX, "".~r2+24(FP)
						//将返回结果存入调用函数栈帧
        0x0018 00024 (test.go:8)        MOVQ    BX, "".~r3+32(FP)
        0x001d 00029 (test.go:8)        RET
        
golang看起来也是借助了栈
```

- Go 有异常类型吗？
```
准确的来讲，go对于error并没有设置类型，大多数的类型都是我们手动的二次封装。

f, err := os.Open("test.txt")
if err != nil {
    log.Fatal(err)
}
```
- 什么是协程（Goroutine）
```
协程其实不太准确，go的goroutine更像是介于用户态和内核态之间的协程，用go的GPM模型来讲，协程就是被调度运行在实际物理线程上的“线程”，对于os来讲，实际的工作是进程，进程上实现了并行，进程下的线程实现逻辑并行（并发），但是由于这种线程的不断切换，需要陷阱中断，寄存器读写值，切换开销较大，协程的出现就是通过软件层面上减少物理上的线程切换，毕竟goroutine才有2kb，是足够开销的。

goroutine的设计上是让大家尽可能的随意启动goroutine，比如http库中的套接字接受内容处理就直接开一个goroutine，但是对于超大型系统，性能开销是一个永恒的话题，并且goroutine也并发那么完美，因此还是要注意比如goroutine泄露这种场景，pprof是一个很好的工具来检测，还有go trace。
```

- 如何高效地拼接字符串
```
首先字符串为什么是只读的，读法大部分都叫做字符串常量，不解释过于底层的东西，大致理解就是字符串的设计就是一个可以被重复利用的常量，比如你声明了一个小明出生地的变量，他的出生地是不可能改变的，因此如果多次复用，或者换成小明最早的出生地这种，都是复用内存中的一个地址。

var str strings.Builder
for i := 0; i < 1000; i++ {
	str.WriteString("a")
}
因为字符串不能复用，所以直接 += 拼接每次都是新建一个字符串，内存开销巨大。
```

- 什么是 rune 类型
```
ASCII码只需要7bit就可以表示所有的英文字母在内的128字符，但是随着发展，需要适配的世界上的文字也越来越多，因此发明了Unicode，为每一个字符分配一个code point（码点），在go等命名为rune，是int32的别称。

Go 语言中，字符串的底层表示是 byte (8 bit) 序列，而非 rune (32 bit) 序列。例如下面的例子中 语 和 言 使用 UTF-8 编码后各占 3 个 byte，因此 len("Go语言") 等于 8，当然我们也可以将字符串转换为 rune 序列。

fmt.Println(len("Go语言")) // 8 这里指的是字节
fmt.Println(len([]rune("Go语言"))) // 4 这里是码点
```

- 如何判断 map 中是否包含某个 key ？
```
if val, ok := dict["en"]; ok {}

最早的map是不能保证thread safe，因此后来也新出了一个sync.Map
大家对go比较诟病的一点就是语法糖太少，甚至很多语法jian，因此go也有很多开源的thread safe map
```

- Go 支持默认参数或可选参数吗？
```
不支持，首先这两个我觉得从编程逻辑的角度很糟糕，他增加了函数的不确定性。
因为这种的参数支持会让函数的命名都产生问题，因为不确定性，你可以说它有很多好处，但是，我一个struct封装你还有什么优势，现在还可以对struct中的字段打各种tag，因此我觉得这俩完全没有存在的必要。

go支持可变参数。
```

- defer的执行顺序
```
后进先出

之前研究过defer的性能开销，在超量调用下defer还是有一定损耗的，貌似现在的版本对这个有了改进，之前的操作大多是在栈上操作的。

还有对于defer用法，很多人喜欢用它来关闭套接字等。我个人认为这种做法不是错的，但是还是有点违背编程的逻辑性。比如我们就打开一个文件，然后读取一下的顺序逻辑，我们加上一个defer就很奇怪，因为它要等这个func完全调用结束才关掉，我们为什么不在读取完后直接关掉呢，那么defer的场景大多不都是关闭资源，这样看来是否有问题，我认为关闭资源，用defer的初衷是避免忘记，因此可以在一个func多个资源开启的时候对其进行同一个的defer关闭，当然看很多开源代码并非这样写的。defer的另一个主要的作用是配合recover来避免goroutine直接down掉整个main。

不少人喜欢面试defer修改的值，我认为真的问题很大，完全就是把defer的设计理念摒弃，只为了用自己的理解来为难求职者。你说defer适合修改值吗？for里用defer合适吗？当然这也属于go的奇技淫巧，对于go熟悉的同志应该不是问题。

func test() {
	i := 0
	defer func() {
		fmt.Println("defer1")
	}()
	defer func() {
		i += 1
		fmt.Println("defer2")
	}()
	return i
}
func main() {
	fmt.Println("return", test())
}
// defer2
// defer1
// return 0
这里i不会改变，说明defer是在return执行之后再执行的，这里的i也不过只是一个临时变量。

func test() (i int) {
	i = 0
	defer func() {
		i += 1
		fmt.Println("defer2")
	}()
	return i
}

func main() {
	fmt.Println("return", test())
}
// defer2
// return 1
这里的i已经不是临时变量了，因此会被defer修改。
```

- 如何交换 2 个变量的值？
```
go的一个语法糖啊。
a, b = b, a
其实实现也不困难，就是在栈上进行几次push pop就ok了
https://www.zhihu.com/question/54500937
```

- Go 语言 tag 的用处？
```
简单来讲就是对于字段的注解吧。

经常使用的比如微服务之间rpc通信，我们都是需要通过字节来传输的，我们将结构体字段打上tag后进行Marshall，这样才能在unmarshall的时候还原出正确的结构体。

还有orm之类的，tag就是让字段的语义更加丰富。

至于如何识别tag，当然是利用反射。
很多同志说golang的反射性能太差，但很多同志却说不怎么用反射，慢就慢吧。科科，其实很多地方都是需要反射的，比如整个tag的解析，go如果在加上generic岂不是会更慢。。。
```

- 如何判断 2 个字符串切片（slice) 是相等的？
```
reflect.DeepEqual(a, b)

或者手动
func StringSliceEqualBCE(a, b []string) bool {
	if len(a) != len(b) {
    	return false
    }
    if (a == nil) != (b == nil) {
    	return false
    }
    b := b[:len(a)]
    for i, v := range a {
    	if v != b[i] {
        	return false
        }
    }
}

return true

因为slice本身吧，其实是一个结构体，每个结构体肯定是对应不同的地址了，所以直接==比较肯定不可以，目前来看对于我们还是要减少对于实际业务中的slice比较，这个场景多用于测试。
```

- 字符串打印时，%v 和 %+v 的区别
```
v 和 %+v 都可以用来打印 struct 的值，区别在于 %v 仅打印各个字段的值，%+v 还会打印各个字段的名称。
type Stu struct {
	Name string
}

func main() {
	fmt.Printf("%v\n", Stu{"Tom"}) // {Tom}
	fmt.Printf("%+v\n", Stu{"Tom"}) // {Name:Tom}
}

但如果结构体定义了 String() 方法，%v 和 %+v 都会调用 String() 覆盖默认值。

实际上对于一个完善的结构体来讲，String()是需要我们设计的，我们要对手动可操作的内容进行比较明确的把控来避免各种意外的产生。
```

- Go 语言中如何表示枚举值(enums)
```
const + iota

type StuType int32

const (
	Type1 StuType = iota
	Type2
	Type3
	Type4
)

func main() {
	fmt.Println(Type1, Type2, Type3, Type4) // 0, 1, 2, 3
}

枚举在Java中并非只是这么简单，实际上对于go这种编程风格来讲，这种enum也还算ok了，对于Java来讲，这压根就算不上enum

https://stackoverflow.com/questions/14426366/what-is-an-idiomatic-way-of-representing-enums-in-go
```
- 空 struct{} 的用途
```
节省内存，其实各种用法初衷和作用也都是基于这一点的
fmt.Println(unsafe.Sizeof(struct{}{})) // 0

map只关注key，不关注value，可以用空struc{}填充

channel并发的信号

还有一个比较特殊的
type Lamp struct{}

func (l Lamp) On() {
        println("On")

}
func (l Lamp) Off() {
        println("Off")
}
这种就是只声明方法，但是没有字段。

```

- init() 函数是什么时候执行的？
```
import -> const -> var -> init -> main
使用init要有个全局把控，不然不仅看起来糟糕，还可能会出现意外的情况。
```

- Go 语言的局部变量分配在栈上还是堆上？
```
这个看编译器！但是人脑还比不过编译器，大部分我们还是能看出的。比如最常见的内存逃逸。
func foo() *int {
	v := 11
	return &v
}

func main() {
	m := foo()
	println(*m) // 11
}

foo() 函数中，如果 v 分配在栈上，foo 函数返回时，&v 就不存在了，但是这段函数是能够正常运行的。Go 编译器发现 v 的引用脱离了 foo 的作用域，会将其分配在堆上。因此，main 函数中仍能够正常访问该值。

还有go的切片初始化，如果超过一定容量也是会直接分配到堆上的。
还有等等，从性能上我们希望尽量使用栈而不是堆，但只是希望，现代编译器已经很聪明了，不太需要在这块考虑太多。比如rust直接就不会让内存逃逸的事情发生。
```

- 2 个 interface 可以比较吗？
```
interface包含两个字段，T,V 也就是type和value。所以可以使用==比较但是要考虑两个是否都相等或者都为nil。

type Stu struct {
	Name string
}

type StuInt interface{}

func main() {
	var stu1, stu2 StuInt = &Stu{"Tom"}, &Stu{"Tom"}
	var stu3, stu4 StuInt = Stu{"Tom"}, Stu{"Tom"}
	fmt.Println(stu1 == stu2) // false
	fmt.Println(stu3 == stu4) // true
}
stu1 和 stu2 对应的类型是 *Stu，值是 Stu 结构体的地址，两个地址不同，因此结果为 false。
stu3 和 stu3 对应的类型是 Stu，值是 Stu 结构体，且各字段相等，因此结果为 true。
```

- 两个 nil 可能不相等吗？
```
可能。
只要记住比较的规则就好了
1.两个接口值比较的时候，先type后value
2.接口值和非接口值比较，先将非接口值转换为接口值，再比较。

	var p *int = nil
	var i interface{} = p
	fmt.Println(i == p) // true
	fmt.Println(p == nil) // true
	fmt.Println(i == nil) // false
    
首先将p赋值给i，这个时候i的type为*int，value为nil
i == p是非接口转接口，相等
p == nil直接比较值
i == nil 先将非接口nil转为接口，这时候type=nil，value=nil，显然和i的type不相等
```

- 简述 Go 语言GC(垃圾回收)的工作原理
```
一次完整的 GC 分为四个阶段：

1）标记准备(Mark Setup，需 STW)，打开写屏障(Write Barrier)
2）使用三色标记法标记（Marking, 并发）
3）标记结束(Mark Termination，需 STW)，关闭写屏障。
4）清理(Sweeping, 并发)

gc不只是背背就能学会的，我也不懂，不会装懂。
```

- 函数返回局部变量的指针是否安全？
```
go会内存逃逸分析，是安全的。
```

- 非接口非接口的任意类型 T() 都能够调用 *T 的方法吗？反过来呢？
```
*T中的T必须是可寻址的。

先列一下哪些不可寻址
	字符串中的字节，这也是一种认为的避免，可寻址的大多是可修改的，因为寻址后的没有地址保护。
    map中的元素，因为map会不停的扩容之类的，但是slice却可以，slice明明也会扩容底层数组。
    常量
    包级别函数
    
    
type T string

func (t *T) hello() {
	fmt.Println("hello")
}

func main() {
	var t1 T = "ABC"
	t1.hello() // hello
	const t2 T = "ABC"
	t2.hello() // error: cannot call pointer method on t
}
```

- 无缓冲的 channel 和 有缓冲的 channel 的区别？
```
channel底层是有一个ring buffer的，里面的recvq和sendq是关键。这个可以参考go夜读里的欧长坤的分享。
https://www.bilibili.com/video/BV1g4411R7p5?from=search&seid=8639757877535805251

对于无缓冲的 channel，发送方将阻塞该信道，直到接收方从该信道接收到数据为止，而接收方也将阻塞该信道，直到发送方将数据发送到该信道中为止。

对于有缓存的 channel，发送方在没有空插槽（缓冲区使用完）的情况下阻塞，而接收方在信道为空的情况下阻塞。

func main() {
	st := time.Now()
	ch := make(chan bool)
	go func ()  {
		time.Sleep(time.Second * 2)
		<-ch
	}()
	ch <- true  // 无缓冲，发送方阻塞直到接收方接收到数据。
	fmt.Printf("cost %.1f s\n", time.Now().Sub(st).Seconds())
	time.Sleep(time.Second * 5)
}

func main() {
	st := time.Now()
	ch := make(chan bool, 2)
	go func ()  {
		time.Sleep(time.Second * 2)
		<-ch
	}()
	ch <- true
	ch <- true // 缓冲区为 2，发送方不阻塞，继续往下执行
	fmt.Printf("cost %.1f s\n", time.Now().Sub(st).Seconds()) // cost 0.0 s
	ch <- true // 缓冲区使用完，发送方阻塞，2s 后接收方接收到数据，释放一个插槽，继续往下执行
	fmt.Printf("cost %.1f s\n", time.Now().Sub(st).Seconds()) // cost 2.0 s
	time.Sleep(time.Second * 5)
}
```

- 什么是协程泄露(Goroutine Leak)？
```
协程泄露其实就是不停的创建goroutine，对于goroutine的设计我们无需手动的关闭，当当前进程结束，goroutine也会消失。

我们只需要关注哪些情况会一直不停创建goroutine。

1.缺少接收器，导致发送阻塞
func query() int {
	ch := make(chan int)
	for i := 0; i < 1000; i++ {
		go func() { ch <- 0 }()
	}
	return <-ch
}

func main() {
	for i := 0; i < 4; i++ {
		query()
		fmt.Printf("goroutines: %d\n", runtime.NumGoroutine())
	}
}
// goroutines: 1001
// goroutines: 2000
// goroutines: 2999
// goroutines: 3998

2.缺少发送器，导致接收阻塞
3.死锁，用go race来进行竞争条件检查，这就涉及个人要求水平和习惯了。
	互斥条件：一个资源只能被一个物理进程使用，也可以认为是一个逻辑线程
    占有等待：一个进程请求资源但是被阻塞会一直等待，并且这个进程占有的资源也不会释放，这就会造成a想要b的b想要a的，也就是比较直观的死锁了
    不能强行占有：进程没办法强行占有另一个进程的资源
    循环依赖关系：a->b->a
4.无限循环创建，我们一般需要利用context，还有一个请求错误次数的监听机制
```

- Go 可以限制运行时操作系统线程的数量吗？
```
runtime.GOMAXPROCS(1) // 限制同时执行Go代码的操作系统线程数为 1

这个什么时候适合手动使用呢，当然是明确的io和cpu计算场景。

默认是1-1的，如果cpu密集的，我们设置1-m（>1）则会造成一个问题，那就是线程切换过多，反而会降低性能。io的话设置大一点会增加吞吐率。
```
- 如何实现select优先级
```
func work(ch1, ch2 <-chan int, stopCh chan struct{})  {
    for {
        select {
        case <-stopCh:
            return
        case job1 := <-ch1:
            println(job1)
        case job2 := <-ch2:
        priority:
            for {
                select {
                case job1 := <-ch1:
                    println(job1)
                default:
                    break priority
                }
            }
            println(job2)
        }
    }
}
```
有点类似单例模式的两次确认，具体做法就是利用default和label特性，这里要注意的是break没法直接跳出select。

k8s cases
```
// kubernetes/pkg/controller/nodelifecycle/scheduler/taint_manager.go 
func (tc *NoExecuteTaintManager) worker(worker int, done func(), stopCh <-chan struct{}) {
    defer done()

    // 当处理具体事件的时候，我们会希望 Node 的更新操作优先于 Pod 的更新
    // 因为 NodeUpdates 与 NoExecuteTaintManager无关应该尽快处理
    // -- 我们不希望用户(或系统)等到PodUpdate队列被耗尽后，才开始从受污染的Node中清除pod。
    for {
        select {
        case <-stopCh:
            return
        case nodeUpdate := <-tc.nodeUpdateChannels[worker]:
            tc.handleNodeUpdate(nodeUpdate)
            tc.nodeUpdateQueue.Done(nodeUpdate)
        case podUpdate := <-tc.podUpdateChannels[worker]:
            // 如果我们发现了一个 Pod 需要更新，我么你需要先清空 Node 队列.
        priority:
            for {
                select {
                case nodeUpdate := <-tc.nodeUpdateChannels[worker]:
                    tc.handleNodeUpdate(nodeUpdate)
                    tc.nodeUpdateQueue.Done(nodeUpdate)
                default:
                    break priority
                }
            }
            // 在 Node 队列清空后我们再处理 podUpdate.
            tc.handlePodUpdate(podUpdate)
            tc.podUpdateQueue.Done(podUpdate)
        }
    }
}
```

NSQ cases
```
for msg := range c.incomingMsgChan {
    select {
    case c.memoryMsgChan <- msg:
    default:
        err := WriteMessageToBackend(&msgBuf, msg, c.backend)
        if err != nil {
            // ... handle errors ...
        }
    }
}
Taking advantage of Go’s select statement allows this functionality to be expressed in just a few lines of code: the default case above only executes if memoryMsgChan is full.
```


## 小结
什么时候，感觉脑子就像泉水一样，就算入门了。所以泉水是怎么生成的？
https://zh.wikipedia.org/wiki/%E6%B3%89

## ref
https://geektutu.com/post/qa-golang-c1.html

