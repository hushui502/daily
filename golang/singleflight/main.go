package main

import (
	"errors"
	"golang.org/x/sync/singleflight"
	"log"
	"sync"
)

var errNotExist = errors.New("not exist")
var g singleflight.Group

func main() {
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			data, err := getData("key")
			if err != nil {
				log.Println(err)
				return
			}
			log.Println(data)
		}()
	}

	wg.Wait()
}

func getData(key string) (string, error) {
	data, err := getDataFromCache(key)
	if err == errNotExist {
		//模拟从db中获取数据
		v, err, _ := g.Do(key, func() (i interface{}, err error) {
			return getDataFromDB(key)
		})
		if err != nil {
			log.Println(err)
			return"", err
		}
		data = v.(string)
		//TOOD: set cache
	} else if err != nil {
		return"", err
	}

	return data, nil
}

func getDataFromCache(key string) (string, error) {
	return "", errNotExist
}

func getDataFromDB(key string) (string, error) {
	log.Printf("get %s from database", key)
	return"data", nil
}


