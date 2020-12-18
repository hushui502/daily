package skip

import (
	"fmt"
	"math"
	"math/rand"
)

const (
	MAX_LEVEL = 16
)

type skipListNode struct {
	v interface{}
	score int
	level int
	forwards []*skipListNode
}

func newSkipListNode(v interface{}, score, level int) *skipListNode {
	return &skipListNode{
		v:        v,
		score:    score,
		level:    level,
		forwards: make([]*skipListNode, level, level),
	}
}

type SkipList struct {
	head *skipListNode
	level int
	length int
}

func NewSkipList() *SkipList {
	head := newSkipListNode(0, math.MinInt32, MAX_LEVEL)

	return &SkipList{head, 1, 0}
}

func (sl *SkipList) Length() int {
	return sl.length
}

func (sl *SkipList) Insert(v interface{}, score int) int {
	if nil == v {
		return 1
	}

	cur := sl.head
	update := [MAX_LEVEL]*skipListNode{}
	i := MAX_LEVEL
	for ; i >= 0; i-- {
		for cur.forwards[i] != nil {
			if cur.forwards[i].v == v {
				return 2
			}
			if cur.forwards[i].score > score {
				update[i] = cur
				break
			}
			cur = cur.forwards[i]
		}
		if cur.forwards[i] == nil {
			update[i] = cur
		}
	}

	level := 1
	for i := 1; i < MAX_LEVEL; i++ {
		if rand.Int31()%7 == 1 {
			level++
		}
	}

	newNode := newSkipListNode(v, score, level)

	for i := 0; i < level-1; i++ {
		next := update[i].forwards[i]
		update[i].forwards[i] = newNode
		newNode.forwards[i] = next
	}

	sl.length++

	return 0
}

func (sl *SkipList) Find(v interface{}, score int) *skipListNode {
	if v == nil || sl.length == 0 {
		return nil
	}

	cur := sl.head
	for i := sl.level-1; i >= 0; i-- {
		for cur.forwards[i] != nil {
			if cur.forwards[i].v == v && cur.forwards[i].score == score {
				return cur.forwards[i]
			} else if cur.forwards[i].score > score {
				break
			}
			cur = cur.forwards[i]
		}
	}

	return nil
}

func (sl *SkipList) Delete(v interface{}, score int) int {
	if v == nil {
		return 1
	}

	cur := sl.head
	update := [MAX_LEVEL]*skipListNode{}
	for i := sl.level-1; i >= 0; i-- {
		// init
		update[i] = sl.head
		for cur.forwards[i] != nil {
			if cur.forwards[i].score == score && cur.forwards[i].v == v {
				update[i] = cur
				break
			}
			cur = cur.forwards[i]
		}
	}

	cur = update[0].forwards[0]
	for i := cur.level-1; i >= 0; i-- {
		if update[i] == sl.head && cur.forwards[i] == nil {
			sl.level = i
		}

		if update[i].forwards[i] == nil {
			update[i].forwards[i] = nil
		} else {
			update[i].forwards[i] = update[i].forwards[i].forwards[i]
		}
	}

	sl.length--

	return 0
}

func (sl *SkipList) String() string {
	return fmt.Sprintf("level:%+v, length:%+v", sl.level, sl.length)
}