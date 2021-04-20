package ttlcache

import (
	"sync"
	"time"
)

// Cache is a synchronised map of items that auto-expire once stale
type Cache struct {
	mutex sync.RWMutex
	ttl   time.Duration
	items map[string]*Item
}

// NewCache is a helper to create instance of the Cache struct
func NewCache(duration time.Duration) *Cache {
	cache := &Cache{
		ttl:   duration,
		items: map[string]*Item{},
	}
	cache.startCleanupTimer()
	return cache
}

// Set is a thread-safe way to add new items to the map
func (cache *Cache) Set(key string, value string) {
	cache.mutex.Lock()
	item := &Item{data: value}
	item.touch(cache.ttl)
	cache.items[key] = item
	cache.mutex.Unlock()
}

func (cache *Cache) Get(key string) (data string, found bool) {
	cache.mutex.RLock()
	item, exists := cache.items[key]
	if !exists || item.expired() {
		data = ""
		found = false
	} else {
		item.touch(cache.ttl)
		data = item.data
		found = true
	}

	cache.mutex.RUnlock()
	return
}

func (cache *Cache) cleanUp() {
	cache.mutex.Lock()
	for key, item := range cache.items {
		if item.expired() {
			delete(cache.items, key)
		}
	}
	cache.mutex.Unlock()
}

func (cache *Cache) Count() int {
	cache.mutex.RLock()
	count := len(cache.items)
	cache.mutex.RUnlock()
	return count
}

func (cache *Cache) startCleanupTimer() {
	duration := cache.ttl
	if duration < time.Second {
		duration = time.Second
	}
	ticker := time.Tick(duration)
	go func() {
		for {
			select {
			case <-ticker:
				cache.cleanUp()
			default:
			}
		}
	}()
}
