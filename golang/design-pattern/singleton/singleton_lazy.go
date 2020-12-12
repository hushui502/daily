package singleton

import "sync"

var (
	lazySingleton *Singleton
	once = &sync.Once{}
)

func GetLazySingelton() *Singleton {
	if lazySingleton == nil {
		// 借助golang的once只会对该func调用一次，不必像java那样多次判断是否为空
		once.Do(func() {
			lazySingleton = &Singleton{}
		})
	}

	return lazySingleton
}
