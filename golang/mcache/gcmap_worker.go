package mcache

import (
	"context"
	"sync"
	"time"
)

var (
	gcInstance *GC
	loadInstance = false
)

type keyset struct {
	keys [2][]string
	cur int
	mutex sync.Mutex
}

func (kset *keyset) len() int {
	return len(kset.keys[kset.cur])
}

func (kset *keyset) append(key string) {
	kset.mutex.Lock()
	defer kset.mutex.Unlock()

	kset.keys[kset.cur] = append(kset.keys[kset.cur], key)
}

func (kset *keyset) swap() []string {
	kset.mutex.Lock()
	defer kset.mutex.Unlock()

	keys := kset.keys[kset.cur]
	kset.keys[kset.cur] = kset.keys[kset.cur][:0]

	// the value of kset.cur only between 0-1
	// swap 1->0 or 0->1
	// 0+1 & 1 ==> 1
	// 1+1 & 1 ==> 0
	kset.cur = (kset.cur + 1) & 0x1

	return keys
}

type GC struct {
	storage SafeMap
	keyChan chan string
}

func NewGC(ctx context.Context, store SafeMap) *GC {
	if loadInstance {
		return gcInstance
	}

	gc := new(GC)
	gc.storage = store
	gc.keyChan = make(chan string, 10000)
	go gc.ExpireKey(ctx)

	gcInstance = gc
	loadInstance = true

	return gc
}

func (gc GC) LenBufferKeyChan() int {
	return len(gc.keyChan)
}

func (gc GC) ExpireKey(ctx context.Context) {
	kset := &keyset{cur:0}
	kset.keys[0] = make([]string, 0, 100)
	kset.keys[1] = make([]string, 0, 100)

	go gc.heartBeatGc(ctx, kset)

	for {
		select {
		case key := <-gc.keyChan:
			kset.append(key)
		case <-ctx.Done():
			loadInstance = false
			return
		}
	}
}

func (gc GC) heartBeatGc(ctx context.Context, kset *keyset) {
	ticker := time.NewTicker(time.Second * 3)
	for {
		select {
		case <-ticker.C:
			if kset.len() == 0 {
				continue
			}
			keys := kset.swap()
			gc.storage.Flush(keys)
		case <-ctx.Done():
			return
		}
	}
}

func (gc GC) Expired(ctx context.Context, key string, duration time.Duration) {
	select {
	case <-time.After(duration):
		gc.keyChan <- key
		return
	case <-ctx.Done():
		return
	}
}