package main

import (
	"fmt"
	"github.com/moby/moby/pkg/pubsub"
	"strings"
	"time"
)

func main() {
	p := pubsub.NewPublisher(100 * time.Microsecond, 10)

	golang := p.SubscribeTopic(func(v interface{}) bool {
		if key, ok := v.(string); ok {
			if strings.HasPrefix(key, "golang") {
				return true
			}
		}
		return false
	})

	docker := p.SubscribeTopic(func(v interface{}) bool {
		if key, ok := v.(string); ok {
			if strings.HasPrefix(key, "docker") {
				return true
			}
		}
		return false
	})

	go p.Publish("hi")
	go p.Publish("golang:golang-----")
	go p.Publish("docker:docker------")
	time.Sleep(1)

	go func() {
		fmt.Println("golang : " , <-golang)
	}()

	go func() {
		fmt.Println("docker : ", <-docker)
	}()

	<-make(chan bool)
}