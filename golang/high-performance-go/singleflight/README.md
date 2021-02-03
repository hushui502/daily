# sync.singleflight

## 背景
缓存在各种场景大量使用，各种意外情况也极其容易发生，比如缓存穿透等。
![](https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/425dc66074394b9bbbe5d188e032f716~tplv-k3u1fbpfcp-watermark.image)

之所以用缓存，目的就是避免请求直接访问持久层数据库，避免给数据库过大的压力。

但是缓存穿透等是不可完全避免的，处理思路应该转到如何处理经过缓存又继续访问数据库的
请求，解决思路就是，再加一层过滤，直接改成”单线程“

![](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/6daa50db7e9b429cae0223cb29c7d6ca~tplv-k3u1fbpfcp-watermark.image)


singleflight看起来十分贴合我们的需求，但也有优缺点。

## singleflight
利用阻塞机制，可以同时访问的只有一个goroutine，因为map是线程不安全的，所以用了很多锁和信号。

### 示例
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