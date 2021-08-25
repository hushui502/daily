# Goroutine 些许知识点

Never start a goroutine without knowing how it will stop.

## goroutine阻塞
主要有两个原因：
- 超时控制
- 流程控制

### Context
一般来讲，开启goroutine处理事务，对于这个事务都有一个完成时间预期
- RPC调用：最大超时时间不会超过用户等待的时间
- 定时任务：执行一次的时间不会超过启动的间隔（定时）

Go中的Context就是对goroutine的生命周期管理
- Cancellation via context.WithCancel
- Timeout via context.WithDeadline

```cassandraql
	req, err := http.NewRequest("GET", "http://www.baidu.com", nil)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(req.Context(), 2 * time.Second)
	defer cancel()

	req = req.WithContext(ctx)
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	fmt.Printf("%v\n", resp.StatusCode)
```

### Channel & Select
流程控制对于其他语言一般来讲就是
```cassandraql
if else
for loop
```
在Go中，结合其特性channel和select也是流程控制的一个重要部分

不要共享变量，用通信！channel的主要作用就是将数据从一个goroutine传递到
另外一个goroutine。

closed channel 永远不会阻塞, close channel不可以再次接收，但是可以无限发送，只不过值为0
```cassandraql
 ch := make(chan bool, 2)
    ch <- true
    ch <- true
    close(ch)

    for v := range ch {
        fmt.Println(v) // called twice
    }
```
nil channel永远阻塞
```cassandraql
var ch chan bool
ch <- true // 
```

buffered/unbuffered channel 介于两者之间，因为会阻塞但是如果”消费掉“会再次不阻塞。

如何确定是否可以读写？select
```cassandraql
select{
	case channel_send_or_receive:
		//Dosomething
	case channel_send_or_receive:
		//Dosomething
	case timeout:
		//...
	case signal:
		// ...
	default:
		//Dosomething
	}
```

## goroutine退出
- 超时退出(如果是处理上下游业务，这种超时是必要的)
- 根据channel可读状态返回(更加符合go的设计理念，但是需要也很容易犯错)

```cassandraql
// 方式一：遍历关闭的 channel
for x := range closedCh {
    fmt.Printf("Process %d\n", x)
}
// 方式二：Select 可读 channel
for {
    select {
        case <-stopCh:
            fmt.Println("Recv stop signal")
            return
        case <- time.After(Duration)
            // ...
        case ...
        case <-t.C:
            fmt.Println("Working .")
    }
}
```

## 完美退出
协程仅仅直接退出，是不够完美的。

至少要包含一下三点：
- 通知协程退出
- 通知确认，协程退出
- 获取协程最终返回的错误
```cassandraql
func (g *Group) Wait() error {
    g.wg.Wait()
    if g.cancel != nil {
        g.cancel()
    }
    return g.err
}
```

## goroutine 三个关键
- 让调用者决定是否启动goroutine，而不是让被调用者内部启动，防止不可知泄露等
- 知道整个goroutine的生命周期，也就是必须要知道开启的goroutine什么时候结束
- 超时控制，管控goroutine，能随意控制goroutine的退出