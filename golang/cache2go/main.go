package main

import (
	"fmt"
	"github.com/muesli/cache2go"
	"time"
)

type myStruct struct {
	text string
	moreData []byte
}

func main() {
	cache := cache2go.Cache("mycache")

	val := myStruct{"this is a test struct", []byte{}}
	cache.Add("someKey", 5*time.Second, &val)
	res, err := cache.Value("someKey")
	if err != nil {

	} else {
		fmt.Println("value is ", res.Data().(*myStruct).text)
	}

	//time.Sleep(6 * time.Second)
	res, err = cache.Value("someKey")
	if err != nil {
		fmt.Println("Item is not cached (anymore).")
	}

	cache.SetAboutToDeleteItemCallback(func(c *cache2go.CacheItem) {
		fmt.Println("Deleting ", c.Key(), c.Data().(*myStruct).text, c.CreatedOn())
	})

	cache.Delete("someKey")

	cache.Flush()
}
