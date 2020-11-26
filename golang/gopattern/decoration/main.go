package main

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

// EXAMPLE 1
func decorator(f func(s string)) func(s string) {
	return func(s string) {
		fmt.Println("started")
		f(s)
		fmt.Println("done")
	}
}

func Hello(s string) {
	fmt.Println(s)
}

// EXAMPLE 2
func Sum1(start, end int64) int64 {
	var sum int64
	sum = 0
	if start > end {
		start, end = end, start
	}
	for i := start; i < end; i++ {
		sum += i
	}

	return sum
}

type SumFunc func(int64, int64) int64

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func timedSumFunc(f SumFunc) SumFunc {
	return func(start int64, end int64) int64 {
		defer func(t time.Time) {
			fmt.Printf("---- Time Elapsed (%s): %v ---\n", getFunctionName(f), time.Since(t))
		}(time.Now())

		return f(start, end)
	}
}

func main() {
	// ex 1
	decorator(Hello)("HELLO WORLD")

	// ex 2
	sum := timedSumFunc(Sum1)
	fmt.Printf("%d \n", sum(1, 1000000))
}


