package main

import (
	"log"
	"ratelimit/leakybucket"
	"ratelimit/simpleratelimit"
	"time"
)

func main() {
	rl := simpleratelimit.New(10, time.Second)
	for i := 0; i < 100; i++ {
		log.Printf("limit result: %v\n", rl.Limit())
	}
	log.Printf("limit result: %v\n", rl.Limit())

	lb := leakybucket.New()
	b, err := lb.Create("leaky_bucket", 10, time.Second)
	if err != nil {
		log.Println(err)
	}
	log.Printf("bucket capacity:%v", b.Capacity())
}
