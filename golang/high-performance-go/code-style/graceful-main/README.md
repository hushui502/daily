## 优雅的main函数

闲着无聊，内容十分初级。

### main函数都要优雅吗？
其实目的不是真的如何写好main，因为现在很多的cli库，让main的作用越来越单一了，就是程序的启动。
这里主要是讲一些go编程的技巧。

### 初级main
```bigquery
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", hello)
	http.ListenAndServe(":8080", mux)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello")
}
```
#### 退出main步骤
- 找出当前的程序pid
- kill -9 {pid}
#### 这种初级main的问题
代码panic后整个程序就挂掉了，即使是goroutine，这也是go的一种fast-fail的设计理念。
问题是我们可能不想让程序这么容易挂掉，毕竟生产代码不是我们的hello world。

### 优雅来了-signal
```bigquery
func main() {
	// create a sig channel, capture syscall signal to this channel
	sig := make(chan os.Signal, 1)
	stopCh := make(chan struct{})
	finishCh := make(chan struct{})
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)

	go func(stopCh, finishCh chan struct{}) {
		for {
			select {
			case <-stopCh:
				fmt.Println("stopped")
				finishCh <- struct{}{}
				return
			default:
				time.Sleep(time.Second)
			}
		}
	}(stopCh, finishCh)

	<-sig
	stopCh <- struct{}{}
	<-finishCh
	fmt.Println("finished")
}
```

### 优雅退出-chan chan
这是Rob Pike给出的一个妙招
```bigquery
func main() {
	sig := make(chan os.Signal)
	stopCh := make(chan chan struct{})
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)

	go func(stopCh chan chan struct{}) {
		for {
			select {
			// channel是进行值拷贝的
			case ch := <-stopCh:
				fmt.Println("stopped")
				ch <- struct{}{}
				return
			default:
				time.Sleep(time.Second)
			}
		}
	}(stopCh)

	<-sig		    
	ch := make(chan struct{})
	stopCh <- ch
	<-ch
	fmt.Println("finished")
}
```
// channel是进行值拷贝的
// channel是引用类型

### 优雅退出-context
```bigquery
func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)
	ctx, cancel := context.WithCancel(context.Background())
	finishedCh := make(chan struct{})

	go func(ctx context.Context, finishedCh chan struct{}) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("stopped")
				finishedCh <- struct{}{}
				return
			default:
				time.Sleep(time.Second)
			}
		}
	}(ctx, finishedCh)

	<-sig
	cancel()
	<-finishedCh
	fmt.Println("finished")
}
```

### 优雅退出-context+waitgroup
```bigquery
func main() {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)
	ctx, cancel := context.WithCancel(context.Background())
	num := 10

	// 用wg来控制多个子goroutine的生命周期
	wg := sync.WaitGroup{}
	wg.Add(num)

	for i := 0; i < num; i++ {
		go func(ctx context.Context) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					fmt.Println("stopped")
					return
				default:
					time.Sleep(time.Duration(i) * time.Second)
				}
			}
		}(ctx)
	}

	<-sig
	cancel()
	// 等待所有的子goroutine都优雅退出
	wg.Wait()
	fmt.Println("finished")
}
```

退出很简单，稍微需要注意的是粒度，这里只是讲了main这块的退出，它不应该控制http handler这种粒度。

### 小结
子goroutine和主goroutine的问题，既然都重要为什么不都在主goroutine里？

Davey Cheney曾讲过尽量不要用goroutine和channel，除非你非常清楚你在干什么。

- 调用链路
- 消息解耦，跨服务使用mq等
- panic-recover
