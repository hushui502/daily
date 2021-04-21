## 为什么要用Context

当一个server在执行，请求来的时候我们启动一个goroutine去处理，然后在这个goroutine当中有对下游服务的rpc调用，也会去请求数据库获取一些数据，这个时候如果下游的处理很慢，但是没有挂，则会导致整个响应时间过长。


![image.png](https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/1d4138e5620a4fbfbf16e89db4ec1a82~tplv-k3u1fbpfcp-watermark.image)

如下图，假如上面的 rpc goroutine 很早就报错了，但是 下面的 db goroutine 又执行了很久，我们最后要返回错误信息，很明显后面 db goroutine 执行的这段时间都是在白白的浪费用户的时间。


![image.png](https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/efe39dd00db040368ca5a1761b0b646f~tplv-k3u1fbpfcp-watermark.image)


context 主要就是用来在多个 goroutine 中设置截止日期、同步信号，传递请求相关值。
每一次 context 都会从顶层一层一层的传递到下面一层的 goroutine 当上面的 context 取消的时候，下面所有的 context 也会随之取消。


![image.png](https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/b0074fc7ad234c49906f394cdedf3c48~tplv-k3u1fbpfcp-watermark.image)


![image.png](https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/9a0c5d9e68d2456fb1aa419b58e0ccea~tplv-k3u1fbpfcp-watermark.image)

## Context
### 使用准则
- 对于server而言，传入请求应该创建一个context
- 通过WithCancel，WithDeadline，WithTimeout创建的Context会同时返回一个cancel方法，这个方法必须被执行，不然会导致 context 泄漏，这个可以通过执行 go vet 命令进行检查
    - 应该将 context.Context 作为函数的第一个参数进行传递，参数命名一般为 ctx 不应该将 Context 作为字段放在结构体中。
- 不要给 context 传递 nil，如果你不知道应该传什么的时候就传递 context.TODO()
    - 不要将函数的可选参数放在 context 当中，context 中一般只放一些全局通用的 metadata 数据，例如 tracing id 等等。
    - context 是并发安全的可以在多个 goroutine 中并发调用

### 函数签名

    ```js
    // 创建一个带有新的 Done channel 的 context，并且返回一个取消的方法
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
    // 创建一个具有截止时间的 context
    // 截止时间是 d 和 parent(如果有截止时间的话) 的截止时间中更早的那一个
    // 当 parent 执行完毕，或 cancel 被调用 或者 截止时间到了的时候，这个 context done 掉
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc)
    // 其实就是调用的 WithDeadline
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
    type CancelFunc
    type Context
    // 一般用于创建 root context，这个 context 永远也不会被取消，或者是 done
    func Background() Context
    // 底层和 Background 一致，但是含义不同，当不清楚用什么的时候或者是还没准备好的时候可以用它
    func TODO() Context
    // 为 context 附加值
    // key 应该具有可比性，一般不应该是 string int 这种默认类型，应该自己创建一个类型
    // 避免出现冲突，一般 key 不应该导出，如果要导出的话应该是一个接口或者是指针
    func WithValue(parent Context, key, val interface{}) Context
    ```

### 源码分析
#### context.Context接口

    ```js
    type Context interface {
        // 返回当前 context 的结束时间，如果 ok = false 说明当前 context 没有设置结束时间
        Deadline() (deadline time.Time, ok bool)
            // 返回一个 channel，用于判断 context 是否结束，多次调用同一个 context done 方法会返回相同的 channel
            Done() <-chan struct{}
        // 当 context 结束时才会返回错误，有两种情况
        // context 被主动调用 cancel 方法取消：Canceled
        // context 超时取消: DeadlineExceeded
        Err() error
            // 用于返回 context 中保存的值, 如何查找，这个后面会讲到
            Value(key interface{}) interface{}
    }
```

#### 默认上下文：context.Background

    ```js
    var (
            background = new(emptyCtx)
            todo       = new(emptyCtx)
        )

    func Background() Context {
        return background
    }

func TODO() Context {
    return todo
}
```

```js
type emptyCtx int

func (*emptyCtx) Deadline() (deadline time.Time, ok bool) {
    return
}

func (*emptyCtx) Done() <-chan struct{} {
    return nil
}

func (*emptyCtx) Err() error {
    return nil
}

func (*emptyCtx) Value(key interface{}) interface{} {
    return nil
}
```

#### 如何取消Context：WithCancel

```js
func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {
    if parent == nil {
        panic("cannot create context from nil parent")
    }
    // 包装出新的 cancelContext
c := newCancelCtx(parent)
       // 构建父子上下文的联系，确保当父 Context 取消的时候，子 Context 也会被取消
       propagateCancel(parent, &c)
       return &c, func() { c.cancel(true, Canceled) }
}
```

简单来看就是绑定到一个parent上
```js
func propagateCancel(parent Context, child canceler) {
    // 首先判断 parent 能不能被取消
done := parent.Done()
          if done == nil {
              return // parent is never canceled
          }

      // 如果可以，看一下 parent 是不是已经被取消了，已经被取消的情况下直接取消 子 context
      select {
          case <-done:
              // parent is already canceled
              child.cancel(false, parent.Err())
                  return
          default:
      }

      // 这里是向上查找可以被取消的 parent context
      if p, ok := parentCancelCtx(parent); ok {
          // 如果找到了并且没有被取消的话就把这个子 context 挂载到这个 parent context 上
          // 这样只要 parent context 取消了子 context 也会跟着被取消
          p.mu.Lock()
              if p.err != nil {
                  // parent has already been canceled
                  child.cancel(false, p.err)
              } else {
                  if p.children == nil {
                      p.children = make(map[canceler]struct{})
                  }
                  p.children[child] = struct{}{}
              }
          p.mu.Unlock()
      } else {
          // 如果没有找到的话就会启动一个 goroutine 去监听 parent context 的取消 channel
          // 收到取消信号之后再去调用 子 context 的 cancel 方法
          go func() {
              select {
                  case <-parent.Done():
                      child.cancel(false, parent.Err())
                  case <-child.Done():
              }
          }()
      }
}
```

接下来我们就看看 cancelCtx 长啥样

```js
type cancelCtx struct {
    Context // 这里保存的是父 Context

        mu       sync.Mutex            // 互斥锁
        done     chan struct{}         // 关闭信号
    children map[canceler]struct{} // 保存所有的子 context，当取消的时候会被设置为 nil
    err      error
}
```
在 Done 方法这里采用了 懒汉式加载的方式，第一次调用的时候才会去创建这个 channel

```js
func (c *cancelCtx) Done() <-chan struct{} {
    c.mu.Lock()
        if c.done == nil {
            c.done = make(chan struct{})
        }
d := c.done
       c.mu.Unlock()
       return d
}
```

Value 方法很有意思，这里相当于是内部 cancelCtxKey 这个变量的地址作为了一个特殊的 key，当查询这个 key 的时候就会返回当前 context 如果不是这个 key 就会向上递归的去调用 parent context 的 Value 方法查找有没有对应的值。

```js
func (c *cancelCtx) Value(key interface{}) interface{} {
    if key == &cancelCtxKey {
        return c
    }
    return c.Context.Value(key)
}
```

Value 方法很有意思，这里相当于是内部 cancelCtxKey 这个变量的地址作为了一个特殊的 key，当查询这个 key 的时候就会返回当前 context 如果不是这个 key 就会向上递归的去调用 parent context 的 Value 方法查找有没有对应的值

```js
func (c *cancelCtx) Value(key interface{}) interface{} {
    if key == &cancelCtxKey {
        return c
    }
    return c.Context.Value(key)
}
```

在前面讲到构建父子上下文之间的关系的时候，有一个去查找可以被取消的父 context 的方法 parentCancelCtx 就用到了这个特殊 value

```js
func parentCancelCtx(parent Context) (*cancelCtx, bool) {
    // 这里先判断传入的 parent 是不是永远不可取消的，如果是就直接返回了
done := parent.Done()
          if done == closedchan || done == nil {
              return nil, false
          }

      // 这里利用了 context.Value 不断向上查询值的特点，只要出现第一个可以取消的 context 的时候就会返回
      // 如果没有的话，这时候 ok 就会等于 false
      p, ok := parent.Value(&cancelCtxKey).(*cancelCtx)
          if !ok {
              return nil, false
          }
      // 这里去判断返回的 parent 的 channel 和传入的 parent 是不是同一个，是的话就返回这个 parent
      p.mu.Lock()
          ok = p.done == done
          p.mu.Unlock()
          if !ok {
              return nil, false
          }
      return p, true
}
```
接下来我们来看最重要的这个 cancel 方法，cancel 接收两个参数，removeFromParent 用于确认是不是把自己从 parent context 中移除，err 是 ctx.Err() 最后返回的错误信息

```js
func (c *cancelCtx) cancel(removeFromParent bool, err error) {
    if err == nil {
        panic("context: internal error: missing cancel error")
    }
    c.mu.Lock()
        if c.err != nil {
            c.mu.Unlock()
                return // already canceled
        }
    c.err = err
        // 由于 cancel context 的 done 是懒加载的，所以有可能存在还没有初始化的情况
        if c.done == nil {
            c.done = closedchan
        } else {
            close(c.done)
        }
    // 循环的将所有的子 context 取消掉
    for child := range c.children {
        // NOTE: acquiring the child's lock while holding parent's lock.
        child.cancel(false, err)
    }
    // 将所有的子 context 和当前 context 关系解除
    c.children = nil
        c.mu.Unlock()

        // 如果需要将当前 context 从 parent context 移除，就移除掉
        if removeFromParent {
            removeChild(c.Context, c)
        }
}
```

#### 超时自动取消如何实现: WithDeadline, WithTimeout
我们先看看比较常用的 WithTimeout, 可以发现 WithTimeout 其实就是调用了 WithDeadline 然后再传入的参数上用当前时间加上了 timeout 的时间

```js
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
    return WithDeadline(parent, time.Now().Add(timeout))
}
```
再来看一下实现超时的 timerCtx，WithDeadline 我们放到后面一点点

```js
type timerCtx struct {
    cancelCtx // 这里复用了 cancelCtx
        timer *time.Timer // Under cancelCtx.mu.

        deadline time.Time // 这里保存了快到期的时间
}
```
Deadline() 就是返回了结构体中保存的过期时间

```js
func (c *timerCtx) Deadline() (deadline time.Time, ok bool) {
    return c.deadline, true
}
```
cancel 其实就是复用了 cancelCtx 中的取消方法，唯一区别的地方就是在后面加上了对 timer 的判断，如果 timer 没有结束主动结束 timer

```js
func (c *timerCtx) cancel(removeFromParent bool, err error) {
    c.cancelCtx.cancel(false, err)
        if removeFromParent {
            // Remove this timerCtx from its parent cancelCtx's children.
            removeChild(c.cancelCtx.Context, c)
        }
    c.mu.Lock()
        if c.timer != nil {
            c.timer.Stop()
                c.timer = nil
        }
    c.mu.Unlock()
}
```
timerCtx 并没有重新实现 Done() 和 Value 方法，直接复用了 cancelCtx 的相关方法

最后我们再看看这个最重要的 WithDeadline 方法

```js
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
    if parent == nil {
        panic("cannot create context from nil parent")
    }

    // 会先判断 parent context 的过期时间，如果过期时间比当前传入的时间要早的话，就没有必要再设置过期时间了
    // 只需要返回 WithCancel 就可以了，因为在 parent 过期的时候，子 context 也会被取消掉
    if cur, ok := parent.Deadline(); ok && cur.Before(d) {
        // The current deadline is already sooner than the new one.
        return WithCancel(parent)
    }

    // 构造相关结构体
c := &timerCtx{
cancelCtx: newCancelCtx(parent),
               deadline:  d,
   }

   // 和 WithCancel 中的逻辑相同，构建上下文关系
   propagateCancel(parent, c)

       // 判断传入的时间是不是已经过期，如果已经过期了就 cancel 掉然后再返回
       dur := time.Until(d)
       if dur <= 0 {
           c.cancel(true, DeadlineExceeded) // deadline has already passed
               return c, func() { c.cancel(false, Canceled) }
       }
   c.mu.Lock()
       defer c.mu.Unlock()

       // 这里是超时取消的逻辑，启动 timer 时间到了之后就会调用取消方法
       if c.err == nil {
           c.timer = time.AfterFunc(dur, func() {
                   c.cancel(true, DeadlineExceeded)
                   })
       }
   return c, func() { c.cancel(true, Canceled) }
}
```
可以发现超时控制其实就是在复用 cancelCtx 的基础上加上了一个 timer 来做定时取消


#### 如何为 Context 附加一些值: WithValue
WithValue 相对简单一点，主要就是校验了一下 Key 是不是可比较的，然后构造出一个 valueCtx 的结构

```js
func WithValue(parent Context, key, val interface{}) Context {
    if parent == nil {
        panic("cannot create context from nil parent")
    }
    if key == nil {
        panic("nil key")
    }
    if !reflectlite.TypeOf(key).Comparable() {
        panic("key is not comparable")
    }
    return &valueCtx{parent, key, val}
}
```

```js
type valueCtx struct {
    Context
        key, val interface{}
}
```

Value 的查找和之前 cancelCtx 类似，都是先判断当前有没有，没有就向上递归，只是在 cancelCtx 当中 key 是一个固定的 key 而已

```js
func (c *valueCtx) Value(key interface{}) interface{} {
    if c.key == key {
        return c.val
    }
    return c.Context.Value(key)
}
```

### 使用场景
- 超时控制

```js
package main

import (
        "context"
        "fmt"
        "time"
       )

// 模拟一个耗时的操作
func rpc() (string, error) {
    time.Sleep(100 * time.Millisecond)
        return "rpc done", nil
}

type result struct {
    data string
        err  error
}

func handle(ctx context.Context, ms int) {
    ctx, cancel := context.WithTimeout(ctx, time.Duration(ms)*time.Millisecond)
        defer cancel()

        r := make(chan result)
        go func() {
            data, err := rpc()
                r <- result{data: data, err: err}
        }()

    select {
        case <-ctx.Done():
            fmt.Printf("timeout: %d ms, context exit: %+v\n", ms, ctx.Err())
        case res := <-r:
                  fmt.Printf("result: %s, err: %+v\n", res.data, res.err)
    }
}

func main() {
    // 这里模拟接受请求，启动一个协程去发起请求
    for i := 1; i < 5; i++ {
        time.Sleep(1 * time.Second)
            go handle(context.Background(), i*50)
    }

    // for test, hang
    time.Sleep(time.Second)
}
```
- 错误取消

```js
package main

import (
        "context"
        "fmt"
        "sync"
        "time"
       )

func f1(ctx context.Context) error {
    select {
        case <-ctx.Done():
            return fmt.Errorf("f1: %w", ctx.Err())
        case <-time.After(time.Millisecond): // 模拟短时间报错
                return fmt.Errorf("f1 err in 1ms")
    }
}

func f2(ctx context.Context) error {
    select {
        case <-ctx.Done():
            return fmt.Errorf("f2: %w", ctx.Err())
        case <-time.After(time.Hour): // 模拟一个耗时操作
                return nil
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
        defer cancel()

        var wg sync.WaitGroup
        wg.Add(2)
        go func() {
            defer wg.Done()
                if err := f1(ctx); err != nil {
                    fmt.Println(err)
                        cancel()
                }
        }()

    go func() {
        defer wg.Done()
            if err := f2(ctx); err != nil {
                fmt.Println(err)
                    cancel()
            }
    }()

    wg.Wait()
}
```
- 跨goroutine数据同步

```js
const requestIDKey int = 0

func WithRequestID(next http.Handler) http.Handler {
    return http.HandlerFunc(
            func(rw http.ResponseWriter, req *http.Request) {
            // 从 header 中提取 request-id
reqID := req.Header.Get("X-Request-ID")
// 创建 valueCtx。使用自定义的类型，不容易冲突
ctx := context.WithValue(
        req.Context(), requestIDKey, reqID)

// 创建新的请求
req = req.WithContext(ctx)

// 调用 HTTP 处理函数
next.ServeHTTP(rw, req)
}
)
    }

// 获取 request-id
func GetRequestID(ctx context.Context) string {
    return ctx.Value(requestIDKey).(string)
}

func Handle(rw http.ResponseWriter, req *http.Request) {
    // 拿到 reqId，后面可以记录日志等等
reqID := GetRequestID(req.Context())
           ...
}

func main() {
handler := WithRequestID(http.HandlerFunc(Handle))
             http.ListenAndServe("/", handler)
}
```
- 防止goroutine泄露

```js
func main() {
    // gen generates integers in a separate goroutine and
    // sends them to the returned channel.
    // The callers of gen need to cancel the context once
    // they are done consuming generated integers not to leak
    // the internal goroutine started by gen.
gen := func(ctx context.Context) <-chan int {
dst := make(chan int)
         n := 1
         go func() {
             for {
                 select {
                     case <-ctx.Done():
                         return // returning not to leak the goroutine
                     case dst <- n:
                             n++
                 }
             }
         }()
     return dst
     }

     ctx, cancel := context.WithCancel(context.Background())
         defer cancel() // cancel when we are finished consuming integers

         for n := range gen(ctx) {
             fmt.Println(n)
                 if n == 5 {
                     break
                 }
         }
}
```


### 缺点
- 最显著的一个就是 context 引入需要修改函数签名，并且会病毒的式的扩散到每个函数上面，不过这个见仁见智，我看着其实还好
- 某些情况下虽然是可以做到超时返回提高用户体验，但是实际上是不会退出相关 goroutine 的，这时候可能会导致 goroutine 的泄漏，针对这个我们来看一个例子

```js
package main

import (
        "net/http"
        _ "net/http/pprof"
        "time"
       )

func main() {
mux := http.NewServeMux()
         mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
                 // 这里阻塞住，goroutine 不会释放的
                 time.Sleep(1000 * time.Second)
                 rw.Write([]byte("hello"))
                 })
handler := http.TimeoutHandler(mux, time.Millisecond, "xxx")
             go func() {
                 if err := http.ListenAndServe("0.0.0.0:8066", nil); err != nil {
                     panic(err)
                 }
             }()
         http.ListenAndServe(":8080", handler)
}

```
查看数据我们可以发现请求返回后， goroutine 其实并未回收，但是如果不阻塞的话是会立即回收的

```js
goroutine profile: total 29
24 @ 0x103b125 0x106cc9f 0x1374110 0x12b9584 0x12bb4ad 0x12c7fbf 0x106fd01
```

我们来看看它的源码，超时控制主要在 ServeHTTP 中实现，我删掉了部分不关键的数据， 我们可以看到函数内部启动了一个 goroutine 去处理请求逻辑，然后再外面等待，但是这里的问题是，当 context 超时之后 ServeHTTP 这个函数就直接返回了，在这里面启动的这个 goroutine 就没人管了

```js
func (h *timeoutHandler) ServeHTTP(w ResponseWriter, r *Request) {
ctx := h.testContext
         if ctx == nil {
             var cancelCtx context.CancelFunc
                 ctx, cancelCtx = context.WithTimeout(r.Context(), h.dt)
                 defer cancelCtx()
         }
     r = r.WithContext(ctx)
         done := make(chan struct{})
         tw := &timeoutWriter{
w:   w,
         h:   make(Header),
         req: r,
         }
panicChan := make(chan interface{}, 1)
               go func() {
                   defer func() {
                       if p := recover(); p != nil {
                           panicChan <- p
                       }
                   }()
                   h.handler.ServeHTTP(tw, r)
                       close(done)
               }()
           select {
               case p := <-panicChan:
                       panic(p)
               case <-done:
                           // ...
               case <-ctx.Done():
                           // ...
           }
}
```

## 总结
context 是一个优缺点都十分明显的包，这个包目前基本上已经成为了在 go 中做超时控制错误取消的标准做法，但是为了添加超时取消我们需要去修改所有的函数签名，对代码的侵入性比较大，如果之前一直都没有使用后续再添加的话还是会有一些改造成本






