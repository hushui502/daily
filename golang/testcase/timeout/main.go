package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	result := testCall(r.Context())
	io.WriteString(w, result+"\n")
}

func testCall(ctx context.Context) string {
	var ts = time.Duration(3) * time.Second
	select {
	case <-ctx.Done():
		log.Println("to cancel")
		return "ctx done"
	case <-time.After(ts):
		log.Printf("timeout %v", ts)
		return "timeout"
	}
}

// timeout case
func main() {
	// use handler
	//srv := http.Server{
	//	Addr:":8080",
	//	Handler:http.HandlerFunc(handler),
	//}
	//if err := srv.ListenAndServe(); err != nil {
	//	fmt.Println("server is failed")
	//}

	// use timeouthandler
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	muxWithMiddlewares := http.TimeoutHandler(mux, 3 * time.Second, "Timeout!")
	http.ListenAndServe(":8080", muxWithMiddlewares)

	// use server params
	//srv := &http.Server{
	//	Handler:      handlers.LoggingHandler(os.Stdout, mux),
	//	Addr:         "localhost:8080",
	//	WriteTimeout: 15 * time.Second,
	//	ReadTimeout:  15 * time.Second,
	//}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	fmt.Fprintf(w, "hello")
}

