package cache

type CachePoolType int

const (
	Mem = iota
	Redis
)

func NewCachePool(cachePoolType CachePoolType) interface{} {
	switch cachePoolType {
	case Redis:
	default:

	}
}
