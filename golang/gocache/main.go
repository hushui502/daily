package main

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"time"
)

func main() {
	c := cache.New(2 * time.Second, 10 * time.Minute)

	c.Set("foo", "bar", cache.DefaultExpiration)

	c.Set("baz", 43, cache.NoExpiration)

	foo, found := c.Get("foo")
	if found {
		fmt.Println(foo)
	}
	//time.Sleep(2 *time.Second)
	<-time.After(2*time.Second)

	foo2, found := c.Get("foo")
	if found {
		fmt.Println(foo2)
	} else {
		fmt.Println("no found")
	}

}
