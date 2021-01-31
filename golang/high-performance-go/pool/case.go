package main

import (
	"encoding/json"
	"sync"
)

type Student struct {
	Name string
	Age int
	Remark [1024]byte
}

// sync.Pool 的大小是可伸缩的，高负载时会动态扩容，存放在池中的对象如果不活跃了会被自动清理。
var studentPool = sync.Pool{
	New: func() interface{} {
		return new(Student)
	},
}

var buf, _ = json.Marshal(&Student{Name: "Hu", Age: 12})

func unmarshal() {
	// 并发度很高的情况下，会产生大量的临时对象，且都分布在堆上，给GC压力很大
	stu := &Student{}
	json.Unmarshal(buf, stu)
}

func unmarshalWithPool() {
	// Get()从池子中取对象，返回值是interface{}
	stu := studentPool.Get().(*Student)
	json.Unmarshal(buf, stu)
	// Put()返回对象池
	studentPool.Put(stu)
}

