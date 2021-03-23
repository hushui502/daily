package expirable_cache

import (
	"container/list"
	"fmt"
	"github.com/pkg/errors"
	"sync"
	"time"
)

// Cache defines cache interface
type Cache interface {
	fmt.Stringer
	Set(key string, value interface{}, ttl time.Duration)
	Get(key string) (interface{}, bool)
	Peek(key string) (interface{}, bool)
	Keys() []string
	Len() int
	Invalidate(key string)
	InvalidateFn(fn func(key string) bool)
	RemoveOldest()
	DeleteExpired()
	Purge()
	Stat() Stats
}

// Stats provides statics for cache
type Stats struct {
	Hits, Misses int	// cache effectiveness
	Added, Evicted int	// number of added and evicted records
}

// cacheImpl provides Cache interface implementations
type cacheImpl struct {
	ttl time.Duration
	maxKeys int
	isLRU bool
	onEvicted func(key string, value interface{})
	sync.Mutex
	stat Stats
	items map[string]*list.Element
	evictList *list.List
}

type cacheItem struct {
	expiresAt time.Time
	key string
	value interface{}
}

const noEvictionTTL = time.Hour * 24 * 365 * 10

var _ Cache = (*cacheImpl)(nil)

// NewCache returns a new Cache
// Default MaxKeys is unlimited.
// Default TTL is 10 years, expirable cache is 5 minutes.
// Default eviction mode is LRC, appropriate option allow to change it to LRU.
func NewCache(options ...Option) (Cache, error) {
	res := cacheImpl{
		items: map[string]*list.Element{},
		evictList: list.New(),
		ttl: noEvictionTTL,
		maxKeys: 0,
	}

	for _, opt := range options {
		if err := opt(&res); err != nil {
			return nil, errors.Wrap(err, "failed to set cache option")
		}
	}

	return &res, nil
}

func (c *cacheImpl) Set(key string, value interface{}, ttl time.Duration) {
	c.Lock()
	defer c.Unlock()

	now := time.Now()
	if ttl == 0 {
		ttl = c.ttl
	}

	// Checking for existing item
	if ent, ok := c.items[key]; ok {
		c.evictList.MoveToFront(ent)
		ent.Value.(*cacheItem).value = value
		ent.Value.(*cacheItem).expiresAt = now.Add(ttl)
		return
	}

	// Add new item
	ent := &cacheItem{key: key, value: value, expiresAt: now.Add(ttl)}
	entry := c.evictList.PushFront(ent)
	c.items[key] = entry
	c.stat.Added++

	// Remove oldest entry if it is expired
	if c.ttl != noEvictionTTL || ttl != noEvictionTTL {
		c.removeOldestIfExpired()
	}

	// verify size not exceeded
	if c.maxKeys > 0 && len(c.items) > c.maxKeys {
		c.removeOldest()
	}
}

func (c *cacheImpl) Get(key string) (interface{}, bool) {
	c.Lock()
	defer c.Unlock()

	if ent, ok := c.items[key]; ok {
		if time.Now().After(ent.Value.(*cacheItem).expiresAt) {
			c.stat.Misses++
			return nil, false
		}
		if c.isLRU {
			c.evictList.MoveToFront(ent)
		}
		c.stat.Hits++
		return ent.Value.(*cacheItem).value, true
	}
	c.stat.Misses++
	return nil, false
}

func (c *cacheImpl) Peek(key string) (interface{}, bool) {
	c.Lock()
	defer c.Unlock()

	if ent, ok := c.items[key]; ok {
		if time.Now().After(ent.Value.(*cacheItem).expiresAt) {
			c.stat.Misses++
			return nil, false
		}
		c.stat.Hits++
		return ent.Value.(*cacheItem).value, true
	}

	c.stat.Misses++
	return nil, false
}

func (c *cacheImpl) Keys() []string {
	c.Lock()
	defer c.Unlock()

	return c.keys()
}

// Len returns count of items in cache, including expired
func (c *cacheImpl) Len() int {
	c.Lock()
	defer c.Unlock()

	return c.evictList.Len()
}

// invalidate key from the cache
func (c *cacheImpl) Invalidate(key string) {
	c.Lock()
	defer c.Unlock()

	if ent, ok := c.items[key]; ok {
		c.removeElement(ent)
	}
}

func (c *cacheImpl) InvalidateFn(fn func(key string) bool) {
	c.Lock()
	defer c.Unlock()

	for key, ent := range c.items {
		if fn(key) {
			c.removeElement(ent)
		}
	}
}

func (c *cacheImpl) RemoveOldest() {
	c.Lock()
	defer c.Unlock()

	c.removeOldest()
}

func (c *cacheImpl) DeleteExpired() {
	c.Lock()
	defer c.Unlock()

	for _, key := range c.keys() {
		if time.Now().After(c.items[key].Value.(*cacheItem).expiresAt) {
			c.removeElement(c.items[key])
		}
	}
}

func (c *cacheImpl) Purge() {
	c.Lock()
	defer c.Unlock()

	for k, v := range c.items {
		delete(c.items, k)
		c.stat.Evicted++
		if c.onEvicted != nil {
			c.onEvicted(k, v.Value.(*cacheItem).value)
		}
	}
	c.evictList.Init()
}

func (c *cacheImpl) Stat() Stats {
	c.Lock()
	defer c.Unlock()
	return c.stat
}

func (c *cacheImpl) String() string {
	stats := c.Stat()
	size := c.Len()
	return fmt.Sprintf("Size: %d, Stats: %+v (%0.1f%%)", size, stats, 100*float64(stats.Hits)/float64(stats.Hits+stats.Misses))
}

func (c *cacheImpl) keys() []string {
	keys := make([]string, 0, len(c.items))

	for ent := c.evictList.Back(); ent != nil; ent = ent.Prev() {
		keys = append(keys, ent.Value.(*cacheItem).key)
	}

	return keys
}

func (c *cacheImpl) removeOldest() {
	ent := c.evictList.Back()
	if ent != nil {
		c.removeElement(ent)
	}
}

func (c *cacheImpl) removeOldestIfExpired() {
	ent := c.evictList.Back()
	if ent != nil && time.Now().After(ent.Value.(*cacheItem).expiresAt) {
		c.removeElement(ent)
	}
}


func (c *cacheImpl) removeElement(e *list.Element) {
	c.evictList.Remove(e)
	kv := e.Value.(*cacheItem)
	delete(c.items, kv.key)
	c.stat.Evicted++
	if c.onEvicted != nil {
		c.onEvicted(kv.key, kv.value)
	}
}