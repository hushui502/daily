package singleton

type Singleton struct {}

var singleton *Singleton

func init() {
	singleton = &Singleton{}
}

// 饿汉，其实也不绝对，这里需要说明只用获取实例的方法，因为golang没有私有
func GetSingleton() *Singleton {
	return singleton
}
