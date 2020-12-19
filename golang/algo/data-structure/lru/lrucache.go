package lru

const (
	hostbit = uint64(^uint(0)) == ^uint64(0)
	LENGTH  = 100
)

type lruNode struct {
	prev *lruNode
	next *lruNode

	key int
	value int

	hnext *lruNode
}

type LRUCache struct {
	node []lruNode

	head *lruNode
	tail *lruNode

	capacity int
	used int
}

func Constructor(capacity int) LRUCache {
	return LRUCache{
		node:     make([]lruNode, LENGTH),
		head:     nil,
		tail:     nil,
		capacity: capacity,
		used:     0,
	}
}

func (l *LRUCache) Get(key int) int {
	if l.tail == nil {
		return -1
	}

	if tmp := l.searchNode(key); tmp != nil {
		l.moveToTail(tmp)
		return tmp.value
	}

	return -1
}

func (l *LRUCache) Put(key int, value int) {
	if tmp := l.searchNode(key); tmp != nil {
		tmp.value = value
		l.moveToTail(tmp)
		return
	}
	l.addNode(key, value)
	if l.used > l.capacity {
		l.delete()
	}
}

func (l *LRUCache) addNode(key int, value int) {
	newNode := &lruNode{
		key:   key,
		value: value,
	}

	tmp := &l.node[hash(key)]
	newNode.hnext = tmp.hnext
	tmp.hnext = newNode
	l.used++

	if l.tail == nil {
		l.tail, l.head = newNode, newNode
		return
	}

	l.tail.next = newNode
	newNode.prev = l.tail
	l.tail = newNode
}

func (l *LRUCache) delete() {
	if l.head == nil {
		return
	}

	prev := &l.node[hash(l.head.key)]
	tmp := prev.hnext

	for tmp != nil && tmp.key != l.head.key {
		prev = tmp
		tmp = tmp.hnext
	}

	if tmp != nil {
		return
	}

	prev.next = tmp.hnext
	l.head = l.head.next
	l.head.prev = nil
	l.used--
}

func (l *LRUCache) moveToTail(node *lruNode) {
	if l.tail == node {
		return
	}

	if l.head == node {
		l.head = node.next
		l.head.prev = nil
	} else {
		node.next.prev = node.prev
		node.prev.next = node.next
	}

	node.next = nil
	l.tail.next = node
	node.prev = l.tail

	l.tail = node
}

func (l *LRUCache) searchNode(key int) *lruNode {
	if l.tail == nil {
		return nil
	}

	tmp := l.node[hash(key)].hnext
	for tmp != nil {
		if tmp.key == key {
			return tmp
		}

		tmp = tmp.hnext
	}

	return nil
}

func hash(key int) int {
	if hostbit {
		return (key ^ (key >> 32)) & (LENGTH - 1)
	}

	return (key ^ (key >> 16)) & (LENGTH - 1)
}