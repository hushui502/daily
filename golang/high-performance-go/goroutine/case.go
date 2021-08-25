//package main
//
//import (
//	"context"
//	"fmt"
//	"net/http"
//	"time"
//)
//
//func main() {
//	req, err := http.NewRequest("GET", "http://www.baidu.com", nil)
//	if err != nil {
//		return
//	}
//
//	ctx, cancel := context.WithTimeout(req.Context(), 2 * time.Second)
//	defer cancel()
//
//	req = req.WithContext(ctx)
//	client := http.DefaultClient
//	resp, err := client.Do(req)
//	if err != nil {
//		return
//	}
//
//	fmt.Printf("%v\n", resp.StatusCode)
//}

package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func main() {
	var a = func() {
		fmt.Println("a goroutine")
		panic("panic!")
	}
	Go(a)


	time.Sleep(1 * time.Second)


	err := testCustomErr()
	//if errors.Is(err, CustomError{}) {
	//	fmt.Println("is custom err")
	//}
	//if errors.As(err, &CustomError{}) {
	//	fmt.Println("as custom err")
	//}

	fmt.Printf("%+v", err)

	tr := NewTracker()
	go tr.Run()

	_ = tr.Event(context.Background(), "test")
	_ = tr.Event(context.Background(), "test")
	_ = tr.Event(context.Background(), "test")

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5 * time.Second))
	defer cancel()

	tr.Shutdown(ctx)

}

func testCustomErr() CustomError {
	return CustomError{
		msg: "custom error!",
		err: errors.New("test"),
	}
}

func Go(x func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("recover success!")
			}
		}()

		x()
	}()
}

type temporary interface {
	Temp() bool
}

func IsTemp(err error) bool {
	te, ok := err.(temporary)
	return ok && te.Temp()
}

type CustomError struct {
	msg string
	err error
}

func (e CustomError) Error() string {
	return e.msg
}

func (e *CustomError) Unwarp() error {
	return e.err
}

type Tracker struct {
	ch chan string
	stop chan struct{}
}

func NewTracker() *Tracker {
	return &Tracker{
		ch: make(chan string, 10),
	}
}

func (t *Tracker) Event(ctx context.Context, data string) error {
	select {
	case t.ch <- data:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *Tracker) Run() {
	for data := range t.ch {
		time.Sleep(time.Second)
		fmt.Println(data)
	}
	t.stop <- struct{}{}
}

func (t *Tracker) Shutdown(ctx context.Context) {
	close(t.ch)
	select {
	case <-t.stop:
	case <-ctx.Done():

	}
}

