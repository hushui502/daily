# sync.Once

## 使用场景
常用于单例模式，初始化配置，保持与数据库连接。作用类似init，但它是延迟执行。
- init会在package首次加载时候执行，但是初始化后的对象可能未必一定会被马上使用，
会造成一定的内存空间的浪费
- Once延迟执行，并发场景下也是线程安全的。

### sync.Once初始化条件
- 当访问某个变量的时候，进行初始化读写
- 初始化时候，所有读会被阻塞，直到写结束
- 变量仅初始化一次，初始化后的变量完成后驻留在内存中

## 示例
```cassandraql
type Config struct {
	Server string
	Port   int64
}

var (
	once   sync.Once
	config *Config
)

func ReadConfig() *Config {
	once.Do(func() {
		var err error
		config = &Config{Server: os.Getenv("TT_SERVER_URL")}
		config.Port, err = strconv.ParseInt(os.Getenv("TT_PORT"), 10, 0)
		if err != nil {
			config.Port = 8080 // default port
		}
		log.Println("init config")
	})
	return config
}

func main() {
	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			_ = ReadConfig()
			defer wg.Done()
		}()
	}
	wg.Wait()
	//time.Sleep(time.Second)
}
```

## 源码分析
```cassandraql
package sync

import (
    "sync/atomic"
)

type Once struct {
    // done作为一个初始化的标志位
    done uint32
    m    Mutex
}

func (o *Once) Do(f func()) {
	// Note: Here is an incorrect implementation of Do:
	//
	//	if atomic.CompareAndSwapUint32(&o.done, 0, 1) {
	//		f()
	//	}
	//
	// Do guarantees that when it returns, f has finished.
	// This implementation would not implement that guarantee:
	// given two simultaneous calls, the winner of the cas would
	// call f, and the second would return immediately, without
	// waiting for the first's call to f to complete.
	// This is why the slow path falls back to a mutex, and why
	// the atomic.StoreUint32 must be delayed until after f returns.

	if atomic.LoadUint32(&o.done) == 0 {
		// Outlined slow-path to allow inlining of the fast-path.
		o.doSlow(f)
	}
}

func (o *Once) doSlow(f func()) {
    o.m.Lock()
    defer o.m.Unlock()
    if o.done == 0 {
        defer atomic.StoreUint32(&o.done, 1)
        f()
    }
}
```
```bigquery
if atomic.CompareAndSwapUint32(&o.done, 0, 1) {
	f()
}
```
这种方式如果多个goroutine同时操作，则会导致f不止执行一次。

为什么done，放在结构体的第一位？
源码中是有解释的：

```cassandraql
    // done indicates whether the action has been performed.
    // It is first in the struct because it is used in the hot path.
    // The hot path is inlined at every call site.
    // Placing done first allows more compact instructions on some architectures (amd64/x86),
    // and fewer instructions (to calculate offset) on other architectures.
```
热路径就是指的是频繁执行的一些指令，如果能减少这类指令的长度，则编译成机器码指令也会更少，会提高性能。
为什么放在第一个字段就能减少？

因为我们的结构体本身就是一个分配到堆上的对象，会有一个指针指向它，第一个字段就是和结构体的地址是相同的，类似C中的数组指针。
其他字段我们不仅需要结构体指针，还需要加上偏移量，比如m我们就需要加上第一个done的偏移量。

## 总结
- Once保证了传入的函数只会执行一次，可以替代单例模式，常用在文件加载，初始化这些场景下。
- Once是不能复用的，只要执行过，再传入也不会再执行了
- Once.Do()在执行过程中如果f出现了panic，后面也不会执行了


https://geektutu.com/post/hpg-sync-once.html