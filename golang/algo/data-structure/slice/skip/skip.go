package skip

/*
Package skip defines a skiplist datastructure. That is, a data structure
that probabilistically determines relationships between keys. By doing
so, it becomes easier to program than a binary search tree but maintains
similar speeds
 */

//const p = .5
//
//type lockedSource struct {
//	mu sync.Mutex
//	src rand.Source
//}
//
//func (ls *lockedSource) Int63() (n int64) {
//	ls.mu.Lock()
//	n = ls.src.Int63()
//	ls.mu.Unlock()
//
//	return
//}
//
//func (ls *lockedSource) Seed(seed int64) {
//	ls.mu.Lock()
//	ls.src.Seed(seed)
//	ls.mu.Unlock()
//}
//
//var generator = rand.New(&lockedSource{src: rand.NewSource(time.Now().UnixNano())})
//
//func generateLevel(maxLevel uint8) uint8 {
//	var level uint8
//	for level = uint8(1); level < maxLevel-1; level++ {
//		if generator.Float64() >= p {
//			return level
//		}
//	}
//
//	return level
//}
//
//func insertNode(sl *SkipList, index uint64) (*SkipList, *SkipList) {
//	right := &SkipList{}
//	right.maxLevel = sl.maxLevel
//	right.level = sl.level
//	right.cache = make(nodes, sl.maxLevel)
//	right.posCache = make(widths, sl.maxLevel)
//	right.head = newNode(nil, sl.maxLevel)
//
//}
//
//type SkipList struct {
//	maxLevel, level uint8
//	head *node
//	num uint64
//	cache nodes
//	posCache widths
//}
//
//func (sl *SkipList) init(ifc interface{}) {
//	switch ifc.(type) {
//	case uint8:
//		sl.maxLevel = 8
//	case uint16:
//		sl.maxLevel = 16
//	case uint32:
//		sl.maxLevel = 32
//	case uint64, uint:
//		sl.maxLevel = 64
//	}
//
//	sl.cache = make(nodes, sl.maxLevel)
//	sl.posCache = make(widths, sl.maxLevel)
//	sl.head = newNode(nil, sl.maxLevel)
//}

//func (sl *SkipList) search(cmp common.Comparator, update node, widths widths) (*node, uint64) {
//	if {
//
//	}
//}
//
//func (sl *SkipList) resetMaxLevel() {
//	if sl.level < 1 {
//		sl.level = 1
//		return
//	}
//	for sl.head.forward[sl.level-1] == nil && sl.level > 1 {
//		sl.level--
//	}
//}
//
//func (sl *SkipList)  {
//
//}
//
