package main

import "sync"

// if more than one goroutine modify this service map, bad thing will happen.
// in golang, normal map is unsafe in multiple thread
//var services = map[string]string{}

var (
	services = map[string]string{}
	serviceMu sync.Mutex
)

// although defer is a feature of golang to close some source, but i really think we should not use it casually
func RegisterService(name, addr string) {
	serviceMu.Lock()
	defer serviceMu.Unlock()
	services[name] = addr
}

func LookupService(name string) string {
	serviceMu.Lock()
	defer serviceMu.Unlock()
	return services[name]
}


