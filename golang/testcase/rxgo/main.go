package main

import (
	"context"
	"fmt"
	"github.com/reactivex/rxgo/v2"
	"math/rand"
	"time"
)

type User struct {
	Name string `json:"name"`
	Age int `json:"age"`
}

func main() {
	//observable := rxgo.Just(1, 2, 3, errors.New("unknown"), 4)()
	//ch := observable.Observe()
	//for item := range ch {
	//	if item.Error() {
	//		fmt.Println("error: ", item.E)
	//	} else {
	//		fmt.Println(item.V)
	//	}
	//}


	//observable := rxgo.Just(1, 2, 3, errors.New("unknown"), 4)()
	//<-observable.ForEach(func(v interface{}) {
	//	fmt.Println("received: ", v)
	//}, func(err error) {
	//	fmt.Println("error: ", err)
	//}, func() {
	//	fmt.Println("completed")
	//})


	//observable := rxgo.Create([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
	//	next <- rxgo.Of(1)
	//	next <- rxgo.Of(2)
	//	next <- rxgo.Error(errors.New("unknown"))
	//	next <- rxgo.Of(3)
	//}})
	//ch := observable.Observe()
	//for item := range ch {
	//	if item.Error() {
	//		fmt.Println("err: ", item.E)
	//	} else {
	//		fmt.Println("item: ", item.V)
	//	}
	//}

	//ch := make(chan rxgo.Item)
	//go func() {
	//	for i := 0; i < 10; i++ {
	//		ch <- rxgo.Of(i)
	//	}
	//	close(ch)
	//}()
	//observable := rxgo.FromChannel(ch)
	//for itrm := range observable.Observe() {
	//	fmt.Println(itrm.V)
	//}


	//observable := rxgo.Interval(rxgo.WithDuration(5 * time.Second))
	//for item := range observable.Observe() {
	//	fmt.Println(item.V)
	//}


	//observable := rxgo.Range(0, 3)
	//for item := range observable.Observe() {
	//	fmt.Println(item.V)
	//}


	//observable := rxgo.Just(1, 2, 3)().Repeat(
	//	3, rxgo.WithDuration(1*time.Second))
	//for item := range observable.Observe() {
	//	fmt.Println(item.V)
	//}


	//observable := rxgo.Defer([]rxgo.Producer{func(_ context.Context, ch chan<- rxgo.Item) {
	//	for i := 0; i < 3; i++ {
	//		ch <- rxgo.Of(i)
	//	}
	//}})
	//
	//for item := range observable.Observe() {
	//	fmt.Println(item.V)
	//}
	//
	//for item := range observable.Observe() {
	//	fmt.Println(item.V)
	//}


	//observable := rxgo.Just(1, 2, 3)()
	//observable = observable.Map(func(_ context.Context, i interface{}) (interface{}, error) {
	//	return i.(int)*2 + 1, nil
	//}).Map(func(_ context.Context, i interface{}) (interface{}, error) {
	//	return i.(int)*3 + 2, nil
	//})
	//for item := range observable.Observe() {
	//	fmt.Println(item.V)
	//}


	//observable := rxgo.Just(
	//	User{
	//		Name:"hufan",
	//		Age:12,
	//	},
	//	User{
	//		Name: "libai",
	//		Age:  10,
	//	})()
	//observable = observable.Marshal(json.Marshal)
	//for item := range observable.Observe() {
	//	fmt.Println(string(item.V.([]byte)))
	//}


	//observable := rxgo.Just(
	//	`{"name":"dj","age":18}`,
	//	`{"name":"jw","age":20}`,
	//)()
	//observable = observable.Map(func(_ context.Context, i interface{}) (interface{}, error) {
	//	return []byte(i.(string)), nil
	//}).Unmarshal(json.Unmarshal, func() interface{} {
	//	return &User{}
	//})
	//for item := range observable.Observe() {
	//	fmt.Println(item.V)
	//}


	//observable := rxgo.Just(1, 2, 3, 4)()
	//observable = observable.BufferWithCount(3)
	//for item := range observable.Observe() {
	//	fmt.Println(item.V)
	//}


	//ch := make(chan rxgo.Item)
	//go func() {
	//	i := 0
	//	for range time.Tick(time.Second) {
	//		ch <- rxgo.Of(i)
	//		i++
	//	}
	//}()
	//observable := rxgo.FromChannel(ch).BufferWithTime(rxgo.WithDuration(3 * time.Second))
	//for item := range observable.Observe() {
	//	fmt.Println(item.V)
	//}


	//ch := make(chan rxgo.Item, 1)
	//go func() {
	//	i := 0
	//	for range time.Tick(time.Second) {
	//		ch <- rxgo.Of(i)
	//		i++
	//	}
	//}()
	//observable := rxgo.FromChannel(ch).BufferWithTimeOrCount(rxgo.WithDuration(3*time.Second), 2)
	//for item := range observable.Observe() {
	//	fmt.Println(item.V)
	//}


	//count := 3
	//observable := rxgo.Range(0, 10).GroupBy(count, func(item rxgo.Item) int {
	//	return item.V.(int) % count
	//}, rxgo.WithBufferedChannel(10))
	//for subObservable := range observable.Observe() {
	//	fmt.Println("new observable")
	//	for item := range subObservable.V.(rxgo.Observable).Observe() {
	//		fmt.Println(item.V)
	//	}
	//}


	observable := rxgo.Range(1, 100)
	observable = observable.Map(func(_ context.Context, i interface{}) (interface{}, error) {
		time.Sleep(time.Duration(rand.Int31()))
		return i.(int)*2+1, nil
	}, rxgo.WithCPUPool()).Filter(func(i interface{}) bool {
		return i.(int)%2 == 0
	}).Distinct(func(_ context.Context, i interface{}) (i2 interface{}, err error) {
		return i, nil
	}).Skip(2).Take(10)
	for item := range observable.Observe() {
		fmt.Println(item.V)
	}

}

// Just
func add(value int) func(int) int {
	return func(a int) int {
		return a + value
	}
}

