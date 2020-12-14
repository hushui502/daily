package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func setup()  {

}

func main() {
	// init config
	setup()

	// 监视服务的退出
	done := make(chan error, 2)

	// 控制服务的退出，传入同一个stop，做到只要有一个服务退出了那么另外一个服务也会随之退出
	stop := make(chan struct{}, 0)

	// for debug
	go func() {
		done <- pprof(stop)
	}()

	// 主服务
	go func() {
		done <- app(stop)
	}()

	// stoped 用于判断当前的stop状态
	var stoped bool


	// 循环读取done这个channel，因为只要有一个退出我们就可以stop channel
	for i := 0; i < cap(done); i++ {
		if err := <-done; err != nil {
			log.Printf("server exit error: %+v", err)
		}

		if !stoped {
			stoped = true
			close(stop)
		}
	}
}

func app(stop <-chan struct{}) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	return server(mux, ":8000", stop)
}

func pprof(stop <-chan struct{}) error {
	// 模拟服务意外退出，用于验证一个服务退出，其他服务同时退出的场景
	go func() {
		server(http.DefaultServeMux, ":8081", stop)
	}()

	time.Sleep(5 * time.Second)
	return fmt.Errorf("mock pprof exit")
}

func server(handler http.Handler, addr string, stop <-chan struct{}) error {
	s := http.Server{
		Handler: handler,
		Addr:    addr,
	}

	// 这个 goroutine 我们可以控制退出，因为只要 stop 这个 channel close 或者是写入数据，这里就会退出
	// 同时因为调用了 s.Shutdown 调用之后，http 这个函数启动的 http server 也会优雅退出
	go func() {
		<-stop
		log.Printf("server will exiting, addr: %s", addr)
		s.Shutdown(context.Background())
	}()

	return s.ListenAndServe()
}
