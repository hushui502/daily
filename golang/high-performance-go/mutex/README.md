## 源码部分理解

## mutex易错

### lock/unlock 必须成对出现
会导致panic

### mutex不能被copy
```bigquery
type Counter struct {
	sync.Mutex
	count int
}

func main() {
	var c Counter
	c.Lock()
	defer c.Unlock()
	c.count++
	foo(c)
}

func foo(c Counter) {

	c.Lock()
	defer c.Unlock()
	fmt.Println(c.count)
}
```
使用vet工具
```bigquery
go vet main.go

# command-line-arguments
./main.go:18:6: call of foo copies lock value: command-line-arguments.Counter
./main.go:21:12: foo passes lock by value: command-line-arguments.Counter
```

### 重入

虽然官方不支持，但是自己可以使用一些hacker的手段简单支持

- 通过hacker的⽅式获取到goroutine id，记录下获取锁的goroutine id，它可以实现Locker接⼝。
  mutex是不可重入的，可重入的好处就是可以一定程度上减少死锁
```bigquery
type RecursiveMutex struct {
	sync.Mutex
	owner int64		// 当前持有锁的goroutine id
	recursion int32	// 当前持有锁的goroutine i
}

func (m *RecursiveMutex) Lock() {
	gid := GoID()
	// 如果当前持有锁的goroutine就是这次调⽤的goroutine,说明是重⼊
	if atomic.LoadInt64(&m.owner) == gid {
		m.recursion++
		return
	}
	m.Mutex.Lock()
	// 获得锁的goroutine第⼀次调⽤，记录下它的goroutine id,调⽤次数加1
	atomic.StoreInt64(&m.owner, gid)
	m.recursion = 1
}

func (m *RecursiveMutex) Unlock() {
	gid := GoID()
	// ⾮持有锁的goroutine尝试释放锁，错误的使⽤
	if atomic.LoadInt64(&m.owner) != gid {
		panic(fmt.Sprintf("wrong the owner(%d): %d!", m.owner, gid))
	}
	// 调⽤次数减1
	m.recursion--
	if m.recursion != 0 { // 如果这个goroutine还没有完全释放，则直接返回
		return
	}
	// 此goroutine最后⼀次调⽤，需要释放锁
	atomic.StoreInt64(&m.owner, -1)
	m.Mutex.Unlock()
}

func GoID() int64 {
	var buf [64]byte

	n := runtime.Stack(buf[:], false)

	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine"))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}

	return int64(id)
}
```
- 调⽤Lock/Unlock⽅法时，由goroutine提供⼀个token，⽤来标识它⾃⼰，⽽不是我们通过hacker的⽅式获取到goroutine id，但是，这样⼀来，就不满⾜Locker接⼝了。
```bigquery
// Token⽅式的递归锁
type TokenRecursiveMutex struct {
    sync.Mutex
    token int64
    recursion int32
}
// 请求锁，需要传⼊token
func (m *TokenRecursiveMutex) Lock(token int64) {
    if atomic.LoadInt64(&m.token) == token { //如果传⼊的token和持有锁的token⼀致，说明是递归调⽤
        m.recursion++
        return
    }
        
    m.Mutex.Lock() // 传⼊的token不⼀致，说明不是递归调⽤
    // 抢到锁之后记录这个token
    atomic.StoreInt64(&m.token, token)
    m.recursion = 1
}

// 释放锁
func (m *TokenRecursiveMutex) Unlock(token int64) {
    if atomic.LoadInt64(&m.token) != token { // 释放其它token持有的锁
        panic(fmt.Sprintf("wrong the owner(%d): %d!", m.token, token))
    }
    
    m.recursion-- // 当前持有这个锁的token释放锁
    if m.recursion != 0 { // 还没有回退到最初的递归调⽤
        return
    }
    atomic.StoreInt64(&m.token, 0) // 没有递归调⽤了，释放锁
    m.Mutex.Unlock()
}
```

## try lock

