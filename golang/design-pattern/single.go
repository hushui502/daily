package pattern

import "sync"

type Instance struct {

}

var (
	goInstance *Instance
	once sync.Once
)

func GoInstance() *Instance {
	if goInstance == nil {
		once.Do(func() {
			goInstance = &Instance{}
		})
	}
	return goInstance
}

