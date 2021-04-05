package main

import (
	"container/list"
	"sync"
)

type Lru interface {
	Add(interface{}, interface{}) bool
	Get(interface{}) (interface{}, bool)
}

func NewLru(size int) Lru {
	return &cache{
		linkedList: list.New(),
		cache:      make(map[interface{}]*list.Element),
		lock:       &sync.RWMutex{},
		size:       size,
	}
}

type cache struct {
	linkedList *list.List
	size       int
	cache      map[interface{}]*list.Element
	lock       *sync.RWMutex
}

type entry struct {
	key   interface{}
	value interface{}
}

func (c *cache) Get(key interface{}) (interface{}, bool) {

	c.lock.RLock()
	defer c.lock.RUnlock()

	if e, ok := c.cache[key]; ok {
		c.linkedList.MoveToFront(e)
		if e.Value.(*entry).value == nil {
			return nil, false
		}
		return e.Value.(*entry).value, true
	}
	return nil, false
}

func (c *cache) Add(key, value interface{}) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	if e, ok := c.cache[key]; ok {
		c.linkedList.MoveToFront(e)
		e.Value.(*entry).value = value
		return true
	}

	e := &entry{key: key, value: value}
	ent := c.linkedList.PushFront(e)
	c.cache[key] = ent

	verdict := c.linkedList.Len() > c.size
	if verdict {
		ent := c.linkedList.Back()
		if ent != nil {
			c.linkedList.Remove(ent)
			kv := ent.Value.(*entry)
			delete(c.cache, kv.key)
		}
	}

	return verdict
}
