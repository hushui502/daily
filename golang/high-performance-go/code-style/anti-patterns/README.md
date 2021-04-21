## Common anti-patterns in Go


https://deepsourcehq.hashnode.dev/common-anti-patterns-in-go

### 0. 什么是反模式
当编写代码时没有未来的因素做出考虑时，就会出现反模式。反模式最初可能看起来是一个适当的问题解决方案，但是，实际上，随着代码库的扩大，这些反模式会变得模糊不清，并给我们的代码库添加 “技术债务”。
反模式的一个简单例子是，在编写 API 时不考虑 API 的消费者如何使用它。

### 1. 从导出函数 (exported function) 返回未导出类型 (unexported type) 的值

以小写字母开头的名称（结构字段，函数或变量）不会被导出，并且仅在定义它们的包内可见。
```bigquery
// 反模式
type unexportedType string

func ExportedFunc() unexportedType { 
    return unexportedType("some string")
} 

// 推荐
type ExportedType string
func ExportedFunc() ExportedType { 
    return ExportedType("some string")
}
```

### 2. 空白标识符的不必要使用
在各种情况下，将值赋值给空白标识符是不需要，也没有必要的。如果在 for 循环中使用空白标识符，Go 规范中提到：

如果最后一个迭代变量是空白标识符，则 range 子句等效于没有该标识符的同一子句。

```bigquery
// 反模式
for _ = range sequence { 
    run()
} 
x, _ := someMap[key] 
_ = <-ch 

// 推荐
for range something { 
    run()
} 

x := someMap[key] 
<-ch
```

### 3. 使用循环/多次 append 连接两个切片
将多个切片附加到一个切片时，无需遍历切片并一个接一个地附加 (append) 每个元素。相反，使用一个 append 语句执行此操作会更好，更有效率。

```bigquery
// 反模式
for _, v := range sliceTwo { 
    sliceOne = append(sliceOne, v)
}

// 推荐
sliceOne = append(sliceOne, sliceTwo…)
```

### 4. make 调用中的冗余参数

该 make 函数是一个特殊的内置函数，用于分配和初始化 map、slice 或 chan 类型的对象。
为了使用 make 初始化切片，我们必须提供切片的类型、切片的长度以及切片的容量作为参数。
在使用 make 初始化 map 的情况下，我们需要传递 map 的大小作为参数。

make的这些参数已经具有默认值
- 对于channel，缓冲区容量默认为0
- 对于map，分配的大小默认为较小的起始大小
- 对于slice， 如果省略容量，则默认和长度一致

```bigquery
// 反模式
ch = make(chan int, 0)
sl = make([]int, 1, 1)

// 推荐
ch = make(chan int)
sl = make([]int, 1)

// 特殊情况
const c = 0
ch = make(chan int, c) // 不是反模式
```

### 5. 函数中无用的 return
```bigquery
// 没用的return，不推荐
func alwaysPrintFoofoo() { 
    fmt.Println("foofoo") 
    return
} 

// 推荐
func alwaysPrintFoo() { 
    fmt.Println("foofoo")
}
```

### 6. switch 语句中无用的 break 语句
在 Go 中，switch 语句不会自动 fallthrough。在像 C 这样的编程语言中，如果前一个 case 语句块中缺少 break 语句，则执行将进入下一个 case 语句中。
但是，人们发现，fallthrough 的逻辑在 switch-case 中很少使用，并且经常会导致错误。
因此，包括 Go 在内的许多现代编程语言都将 switch-case 的默认逻辑改为不 fallthrough。
```bigquery
// 反模式
switch s {
case 1: 
    fmt.Println("case one") 
    break
case 2: 
    fmt.Println("case two")
}

// 推荐
switch s {
case 1: 
    fmt.Println("case one")
case 2: 
    fmt.Println("case two")
}

// fallthrough
switch 2 {
case 1: 
    fmt.Print("1") 
    fallthrough
case 2: 
    fmt.Print("2") 
    fallthrough
case 3: fmt.Print("3")
}
```

### 7. 不使用辅助函数执行常见任务
对于一组特定的参数，某些函数具有一些特定表达方式，可以用来简化效率，并带来更好的理解/可读性。
```bigquery
// 反模式
wg.Add(1) // ...some code
wg.Add(-1)

// 推荐
wg.Add(1)
// ...some code
wg.Done()
```

### 8. nil 切片上的冗余检查
nil 切片的长度为零。因此，在计算切片的长度之前，无需检查切片是否为 nil 切片。
```bigquery
// 反模式
if x != nil && len(x) != 0 { // do something
}

// 推荐
if len(x) != 0 { // do something
}
```

### 9. 太复杂的函数字面量
可以删除仅调用单个函数且对函数内部的值没有做任何修改的函数字面量，因为它们是多余的。
可以改为在外部函数直接调用被调用的内部函数。
```bigquery
// 反模式
fn := func(x int, y int) int { return add(x, y) }

// 推荐
add(x, y)
```
刷题党表示坑惨了

### 10. 使用仅有一个 case 语句的 select 语句
select 语句使 goroutine 等待多个通信操作。但是，如果只有一个 case 语句，实际上我们不需要使用 select 语句。在这种情况下，使用简单 send 或 receive 操作即可。
如果我们打算在不阻塞地发送或接收操作的情况处理 channel 通信，则建议在 select 中添加一个 default case 以使该 select 语句变为非阻塞状态。
```bigquery
// 反模式
select {
    case x := <-ch: fmt.Println(x)
} 

// 推荐
x := <-ch
fmt.Println(x)
```
使用default
```bigquery
select {
    case x := <-ch: 
        fmt.Println(x)
    default: 
        fmt.Println("default")
}
```

### 11. context.Context 应该是函数的第一个参数
context.Context 应该是第一个参数，一般命名为 ctx.ctx 应该是 Go 代码中很多函数的（非常）常用参数，由于在逻辑上把常用参数放在参数列表的第一个或最后一个比较好。为什么这么说呢？因为它的使用模式统一，可以帮助我们记住包含该参数。
在 Go 中，由于变量可能只是参数列表中的最后一个，因此建议将 context.Context 作为第一个参数。
各种项目，甚至 Node.js 等都有一些约定，比如错误先回调。
因此，context.Context 应该永远是函数的第一个参数，这是一个惯例。
```bigquery
// 反模式
func badPatternFunc(k favContextKey, ctx context.Context) {    
    // do something
}

// 推荐
func goodPatternFunc(ctx context.Context, k favContextKey) {    
    // do something
}
```