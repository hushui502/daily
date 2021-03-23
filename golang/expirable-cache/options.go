package expirable_cache

import "time"

// Option func type
type Option func(lc *cacheImpl) error

// OnEvicted called automatically for automatically and manually deleted entries
func OnEvicted(fn func(key string, value interface{})) Option {
	return func(lc *cacheImpl) error {
		lc.onEvicted= fn
		return nil
	}
}

func MaxKeys(max int) Option {
	return func(lc *cacheImpl) error {
		lc.maxKeys = max
		return nil
	}
}

func TTL(ttl time.Duration) Option {
	return func(lc *cacheImpl) error {
		lc.ttl = ttl
		return nil
	}
}

func LRU() Option {
	return func(lc *cacheImpl) error {
		lc.isLRU = true
		return nil
	}
}