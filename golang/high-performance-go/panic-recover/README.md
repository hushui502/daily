# panic-recover 用例

## 用例1：避免恐慌导致程序崩溃
```bigquery
package main

import (
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalln(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		go ClientHandler(conn)
	}
}

func ClientHandler(c net.Conn) {
	defer func() {
		if e := recover(); e != nil {
			log.Println("capture a panic: ", e)
			log.Println("recover success!")
		}
		c.Close()
	}()

	panic("unknow error")
}
```

## 用例2：自动重启因为恐慌而退出的协程
```bigquery
package main

import (
	"log"
	"time"
)

func shouldNotExit() {
	for {
		time.Sleep(time.Duration(1) * time.Second)
		if time.Now().UnixNano() & 0x3 == 0 {
			panic("unexpected situation")
		}
	}
}

func NeverExit(name string, f func()) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("协程%s崩溃了，准备重启一个", name)
			go NeverExit(name, f)
		}
	}()
	f()
}

func main() {
	log.SetFlags(0)
	go NeverExit("job#A", shouldNotExit)
	go NeverExit("job#B", shouldNotExit)
	select{} // 永久阻塞主线程
}
```

## 用例3：一些recover调用相当于空操作（No-Op）
```bigquery
package main

func main() {
	defer func() {
		defer func() {
			recover() // 空操作
		}()
	}()
	defer func() {
		func() {
			recover() // 空操作
		}()
	}()
	func() {
		defer func() {
			recover() // 空操作
		}()
	}()
	func() {
		defer recover() // 空操作
	}()
	func() {
		recover() // 空操作
	}()
	recover()       // 空操作
	defer recover() // 空操作
	panic("bye")
}
```

## 用例4：使用panic/recover函数调用模拟长程跳转
```bigquery
package main

import "fmt"

func main() {
	n := func () (result int)  {
		defer func() {
			if v := recover(); v != nil {
				if n, ok := v.(int); ok {
					result = n
				}
			}
		}()

		func () {
			func () {
				func () {
					// ...
					panic(123) // 用恐慌来表示成功返回
				}()
				// ...
			}()
		}()
		// ...
		return 0
	}()
	fmt.Println(n) // 123
}
```

## 用例5：使用panic/recover函数调用来减少错误检查代码
```bigquery
func doSomething() (err error) {
	defer func() {
		err = recover()
	}()

	doStep1()
	doStep2()
	doStep3()
	doStep4()
	doStep5()

	return
}

// 在现实中，各个doStepN函数的原型可能不同。
// 每个doStepN函数的行为如下：
// * 如果已经成功，则调用panic(nil)来制造一个恐慌
//   以示不需继续；
// * 如果本步失败，则调用panic(err)来制造一个恐慌
//   以示不需继续；
// * 不制造任何恐慌表示继续下一步。
func doStepN() {
	...
	if err != nil {
		panic(err)
	}
	...
	if done {
		panic(nil)
	}
}
```

## 在下面的情况下，recover函数调用的返回值为nil：
- 传递给相应panic函数调用的实参为nil；
- 当前协程并没有处于恐慌状态；
- recover函数并未直接在一个延迟函数调用中调用。

在任何时刻，一个协程中只有最新产生的恐慌才能够被恢复。 
换句话说，每个recover调用都试图恢复当前协程中最新产生的且尚未恢复的恐慌。 这解释了为什么上例中的第二个recover调用不会起作用。
```bigquery
// 此程序将带着未被恢复的恐慌1而崩溃退出。
package main

func demo() {
	defer func() {
		defer func() {
			recover() // 此调用将恢复恐慌2
		}()

		defer recover() // 空操作

		panic(2)
	}()
	panic(1)
}

func main() {
	demo()
}
```

## 小结
一般来讲，后两种不推荐。panic-recover的目的不是炫技，而是防止整个main thread挂掉。

一个recover调用只有在它的直接外层调用（即recover调用的父调用）是一个延迟调用，并且此延迟调用（即父调用）
的直接外层调用（即recover调用的爷调用）和当前协程中最新产生并且尚未恢复的恐慌相关联时才起作用。
一个有效的recover调用将最新产生并且尚未恢复的恐慌和与此恐慌相关联的函数调用（即爷调用）剥离开来，
并且返回当初传递给产生此恐慌的panic函数调用的参数。