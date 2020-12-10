package mcache

import (
	"context"
	"time"
)

const TTL_FOREVER = time.Hour * 87660

func (mc *CacheDriver) initStore() (context.Context, context.CancelFunc) {
	ctx, finish := context.WithCancel(context.Background())
	mc.storage = NewStorage()
	mc.gc = NewGC(ctx, mc.storage)

	return ctx, finish
}

type CacheDriver struct {
	ctx context.Context
	closeCtx context.CancelFunc
	storage SafeMap
	gc *GC
	instance *CacheDriver
}

func startInstance() *CacheDriver {
	cdriver := new(CacheDriver)
	ctx, finish := cdriver.initStore()
	cdriver.ctx = ctx
	cdriver.closeCtx = finish

	return cdriver
}

func New() *CacheDriver {
	cdriver := new(CacheDriver)
	ctx, finish := cdriver.initStore()
	cdriver.ctx = ctx
	cdriver.closeCtx = finish

	return cdriver
}

func (mc *CacheDriver) Get(key string) (interface{}, bool) {
	data, ok := mc.storage.Find(key)
	if !ok {
		return Item{}.DataLink, false
	}

	entity := data.(Item)
	if entity.IsExpire() {
		return Item{}.DataLink, false
	}

	return entity.DataLink, true
}

func (mc *CacheDriver) Set(key string, value interface{}, ttl time.Duration) error {
	expire := time.Now().Local().Add(ttl)
	if ttl != TTL_FOREVER {
		go mc.gc.Expired(mc.ctx, key, ttl)
	}

	mc.storage.Insert(key, Item{Key: key, Expire: expire, DataLink: value})

	return nil
}

func (mc *CacheDriver) Remove(key string) {
	mc.storage.Delete(key)
}

func (mc *CacheDriver) Truncate() {
	mc.storage.Truncate()
}

func (mc *CacheDriver) Len() int {
	return mc.storage.Len()
}


func (mc *CacheDriver) GCBufferQueue() int {
	return mc.gc.LenBufferKeyChan()
}

func (mc *CacheDriver) Close() map[string]interface{} {
	mc.closeCtx()
	return mc.storage.Close()
}
