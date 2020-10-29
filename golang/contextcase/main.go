package main

import (
	"context"
	"time"
)

type key int

const (
	userIP = iota
	userID
	logID
)

type Result struct {
	order string
	logistics string
	recommend string
}

func api() (result *Result, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 1)
	defer cancel()

	ctx = context.WithValue(ctx, userIP, "127.0.0.1")
	ctx = context.WithValue(ctx, userID, 666888)
	ctx = context.WithValue(ctx, logID, "123456")

	result = &Result{}

	go func() {
		result.order, err = getOrderDetail(ctx)
	}()

	for {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:

		}
		if err != nil {
			return result, err
		}
		if result.order != "" {
			return result, nil
		}
	}

}

func getOrderDetail(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()

	time.Sleep(time.Millisecond * 700)

	uip := ctx.Value(userIP).(string)
	println(uip)

	return handleTimeout(ctx, func() string {
		return "order"
	})
}

func handleTimeout(ctx context.Context, f func() string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	str := make(chan string)
	go func() {
		str <- f()
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case ret := <-str:
		return ret, nil
	}
}


func main() {
	
}
