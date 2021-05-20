```
# sync.singleflight

## 背景
缓存在各种场景大量使用，各种意外情况也极其容易发生，比如缓存穿透等。
![](https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/425dc66074394b9bbbe5d188e032f716~tplv-k3u1fbpfcp-watermark.image)

之所以用缓存，目的就是避免请求直接访问持久层数据库，避免给数据库过大的压力。

但是缓存穿透等是不可完全避免的，处理思路应该转到如何处理经过缓存又继续访问数据库的
请求，解决思路就是，再加一层过滤，直接改成”单线程“

![](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/6daa50db7e9b429cae0223cb29c7d6ca~tplv-k3u1fbpfcp-watermark.image)

如下图所示，可能存在来自桌面端和移动端的用户有 1000 的并发请求，他们都访问的获取文章列表的接口，获取前 20 条信息，如果这个时候我们服务直接去访问 redis 出现 cache miss 那么我们就会去请求 1000 次数据库，这时可能会给数据库带来较大的压力（这里的 1000 只是一个例子，实际上可能远大于这个值）导致我们的服务异常或者超时。
![](https://img2020.cnblogs.com/blog/1459374/202105/1459374-20210520221929010-766253012.png)

这时候就可以使用 singleflight 库了，直译过来就是单飞，这个库的主要作用就是将一组相同的请求合并成一个请求，实际上只会去请求一次，然后对所有的请求返回相同的结果。

singleflight看起来十分贴合我们的需求，但也有优缺点。

## singleflight
利用阻塞机制，可以同时访问的只有一个goroutine，因为map是线程不安全的，所以用了很多锁和信号。

### 函数签名
```
type Group
// Do 执行函数, 对同一个 key 多次调用的时候，在第一次调用没有执行完的时候
// 只会执行一次 fn 其他的调用会阻塞住等待这次调用返回
// v, err 是传入的 fn 的返回值
// shared 表示是否真正执行了 fn 返回的结果，还是返回的共享的结果
func (g *Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool)

	// DoChan 和 Do 类似，只是 DoChan 返回一个 channel，也就是同步与异步的区别
	func (g *Group) DoChan(key string, fn func() (interface{}, error)) <-chan Result

    // Forget 用于通知 Group 删除某个 key 这样后面继续这个 key 的调用的时候就不会在阻塞等待了
	func (g *Group) Forget(key string)
```


### 示例

demo 1
```cassandraql
type Result string

func find(ctx context.Context, query string) (Result, error) {
	return Result(fmt.Sprintf("result for %q", query)), nil
}

func main() {
	var g singleflight.Group
	const n = 5
	wailted := int32(n)
	done := make(chan struct{})
	key := "https://weibo.com/1227368500/H3GIgngon"
	for i := 0; i < n; i++ {
		go func(j int) {
			v, _, shared := g.Do(key, func() (i interface{}, err error) {
				ret, err := find(context.Background(), key)
				return ret, err
			})
			if atomic.AddInt32(&wailted, -1) == 0 {
				close(done)
			}
			fmt.Printf("index: %d, val: %v, shared: %v\n", j, v, shared)
		}(i)
	}

	select {
	case <-done:
	case <-time.After(time.Second):
		fmt.Println("Do hangs!")
	}
}
```

demo 2
```
var count int32
func getArticle(id int) (article string, err error) {
	atomic.AddInt32(&count, 1)
	time.Sleep(time.Duration(count) * time.Millisecond)

	return fmt.Sprintf("article: %d", id), nil
}

func singleFlightGetArticle(sg *singleflight.Group, id int) (string, error) {
	v, err, _ := sg.Do(fmt.Sprintf("%d", id), func() (interface{}, error) {
		return getArticle(id)
	})

	return v.(string), err
}

func main() {
	time.AfterFunc(time.Duration(1)*time.Second, func() {
		atomic.AddInt32(&count, -count)
	})

	var (
		wg sync.WaitGroup
		now = time.Now()
		n = 1000
		sg = &singleflight.Group{}
	)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			// res, _ := getArticle(1)			// 同时发起 1000 次请求，耗时: 1.0022831s
			res, _ := singleFlightGetArticle(sg, 1)		// 同时发起 1000 次请求，耗时: 1.5119ms
			if res != "article: 1" {
				panic("err")
			}
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Printf("同时发起 %d 次请求，耗时: %s", n, time.Since(now))
}

```

### 源码分析
```
type Group struct {
	mu sync.Mutex       // protects m
	m  map[string]*call // lazily initialized
}
```
Group 结构体由一个互斥锁和一个 map 组成，可以看到注释 map 是懒加载的，所以 Group 只要声明就可以使用，不用进行额外的初始化零值就可以直接使用。call 保存了当前调用对应的信息，map 的键就是我们调用 Do 方法传入的 key

```
type call struct {
	wg sync.WaitGroup

	// 函数的返回值，在 wg 返回前只会写入一次
	val interface{}
	err error

	// 使用调用了 Forgot 方法
	forgotten bool

    // 统计调用次数以及返回的 channel
	dups  int
	chans []chan<- Result
}
```

```
func (g *Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool) {
	g.mu.Lock()

    // 前面提到的懒加载
    if g.m == nil {
		g.m = make(map[string]*call)
	}

    // 会先去看 key 是否已经存在
	if c, ok := g.m[key]; ok {
       	// 如果存在就会解锁
		c.dups++
		g.mu.Unlock()

        // 然后等待 WaitGroup 执行完毕，只要一执行完，所有的 wait 都会被唤醒
		c.wg.Wait()

        // 这里区分 panic 错误和 runtime 的错误，避免出现死锁，后面可以看到为什么这么做
		if e, ok := c.err.(*panicError); ok {
			panic(e)
		} else if c.err == errGoexit {
			runtime.Goexit()
		}
		return c.val, c.err, true
	}

    // 如果我们没有找到这个 key 就 new call
	c := new(call)

    // 然后调用 waitgroup 这里只有第一次调用会 add 1，其他的都会调用 wait 阻塞掉
    // 所以这要这次调用返回，所有阻塞的调用都会被唤醒
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

    // 然后我们调用 doCall 去执行
	g.doCall(c, key, fn)
	return c.val, c.err, c.dups > 0
}
```

这个方法的实现有点意思，使用了两个 defer 巧妙的将 runtime 的错误和我们传入 function 的 panic 区别开来避免了由于传入的 function panic 导致的死锁
```
func (g *Group) doCall(c *call, key string, fn func() (interface{}, error)) {
	normalReturn := false
	recovered := false

    // 第一个 defer 检查 runtime 错误
    defer func() {
  	// 如果既没有正常执行完毕，又没有 recover 那就说明需要直接退出了
  	if !normalReturn && !recovered {
  		c.err = errGoexit
  	}

  	c.wg.Done()
  	g.mu.Lock()
  	defer g.mu.Unlock()

         // 如果已经 forgot 过了，就不要重复删除这个 key 了
  	if !c.forgotten {
  		delete(g.m, key)
  	}

  	if e, ok := c.err.(*panicError); ok {
  		// 如果返回的是 panic 错误，为了避免 channel 死锁，我们需要确保这个 panic 无法被恢复
  		if len(c.chans) > 0 {
  			go panic(e)
  			select {} // Keep this goroutine around so that it will appear in the crash dump.
  		} else {
  			panic(e)
  		}
  	} else if c.err == errGoexit {
  		// 已经准备退出了，也就不用做其他操作了
  	} else {
  		// 正常情况下向 channel 写入数据
  		for _, ch := range c.chans {
  			ch <- Result{c.val, c.err, c.dups > 0}
  		}
  	}
  }()

    // 使用一个匿名函数来执行
	func() {
		defer func() {
			if !normalReturn {
                // 如果 panic 了我们就 recover 掉，然后 new 一个 panic 的错误
                // 后面在上层重新 panic
				if r := recover(); r != nil {
					c.err = newPanicError(r)
				}
			}
		}()

		c.val, c.err = fn()

        // 如果 fn 没有 panic 就会执行到这一步，如果 panic 了就不会执行到这一步
        // 所以可以通过这个变量来判断是否 panic 了
		normalReturn = true
	}()

    // 如果 normalReturn 为 false 就表示，我们的 fn panic 了
    // 如果执行到了这一步，也说明我们的 fn  recover 住了，不是直接 runtime exit
	if !normalReturn {
		recovered = true
	}
}
```

DoChan
Do chan 和 Do 类似，其实就是一个是同步等待，一个是异步返回，主要实现上就是，如果调用 DoChan 会给 call.chans 添加一个 channel 这样等第一次调用执行完毕之后就会循环向这些 channel 写入数据
```
func (g *Group) DoChan(key string, fn func() (interface{}, error)) <-chan Result {
	ch := make(chan Result, 1)
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		c.dups++
		c.chans = append(c.chans, ch)
		g.mu.Unlock()
		return ch
	}
	c := &call{chans: []chan<- Result{ch}}
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	go g.doCall(c, key, fn)

	return ch
}
```
forget 用于手动释放某个 key 下次调用就不会阻塞等待了
```
func (g *Group) Forget(key string) {
	g.mu.Lock()
	if c, ok := g.m[key]; ok {
		c.forgotten = true
	}
	delete(g.m, key)
	g.mu.Unlock()
}
```

### 坑
1. 一个阻塞，全员等待
```
func singleflightGetArticle(sg *singleflight.Group, id int) (string, error) {
	v, err, _ := sg.Do(fmt.Sprintf("%d", id), func() (interface{}, error) {
		// 模拟出现问题，hang 住
		select {}
		return getArticle(id)
	})

	return v.(string), err
}
```
再次允许会出现死锁
```
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [select (no cases)]:
```
解决方案：超时控制
```
func singleflightGetArticle(ctx context.Context, sg *singleflight.Group, id int) (string, error) {
	result := sg.DoChan(fmt.Sprintf("%d", id), func() (interface{}, error) {
		// 模拟出现问题，hang 住
		select {}
		return getArticle(id)
	})

	select {
	case r := <-result:
		return r.Val.(string), r.Err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
```

2. 一个出错，全员出错
   这个本身不是什么问题，因为 singleflight 就是这么设计的，但是实际使用的时候 如果我们一次调用要 1s，我们的数据库请求或者是 下游服务可以支撑 10rps 的请求的时候这会导致我们的错误阈提高，因为实际上我们可以一秒内尝试 10 次，但是用了 singleflight 之后只能尝试一次，只要出错这段时间内的所有请求都会受影响

这种情况我们可以启动一个 Goroutine 定时 forget 一下，相当于将 rps 从 1rps 提高到了 10rps

```
go func() {
       time.Sleep(100 * time.Millisecond)
       // logging
       g.Forget(key)
   }()
```


### 问题分析
如果QPS很高，且设置了请求超时，则可能存在的问题:
- 协程猛增。因为是阻塞模型，所以可能会在某一个时刻大量的协程都在等待没法释放
  项目也不止这一个服务需要协程。
- 内存暴涨。
- 大量的超时报错。
- 后续请求耗时增加。调度问题。

根本原因：
- 阻塞读：缺少超时控制，难以fast fail
- 单并发：可以说是一种很不优雅的控制流量的方式，牺牲了成功率。

如何解决：
### 阻塞读
```cassandraql
// DoChan
	for i := 0; i < n; i++ {
		go func(j int) {
			ch := g.DoChan(key, func() (i interface{}, err error) {
				ret, err := find(context.Background(), key)
				return ret, err
			})

			// create timeout
			timeout := time.After(time.Second)
			var ret singleflight.Result
			select {
			case <-timeout:
				fmt.Println("Timeout!")
				return
			case ret = <-ch:
				fmt.Printf("Index: %d, val: %d, sharead: %d\n", j, ret, ret.Val)
			}
		}(i)
	}
```

### 单并发
```cassandraql
v, _, shared := g.Do(key, func() (interface{}, error) {
    go func() {
        time.Sleep(10 * time.Millisecond)
        fmt.Printf("Deleting key: %v\n", key)
        g.Forget(key)
    }()
    ret, err := find(context.Background(), key)
    return ret, err
})
```

当某个请求执行了1ms还没有结束，则不会阻塞其他的请求。也就是说从以前的最大100QPS更改为大于100QPS

## 总结
singleflight并非银弹，尤其是在可用性极高的场景。

其他方案：
- 放弃这种同步请求的方式，牺牲数据的实时性
- 缓存 准实时的数据 + 异步更新 数据到缓存

https://www.cyningsun.com/01-11-2021/golang-concurrency-singleflight.html
```