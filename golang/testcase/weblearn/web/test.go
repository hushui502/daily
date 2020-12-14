package main

//import (
//	"net/http"
//	"time"
//)
//
//func hello(w http.ResponseWriter, r *http.Request) {
//	w.Write([]byte("hello"))
//
//}
//
//func timeMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		timeStart := time.Now()
//		next.ServeHTTP(w, r)
//		timeElapsed := time.Since(timeStart)
//		println(timeElapsed)
//	})
//}
//
//func timeOut(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		println("timeout start")
//		//next.ServeHTTP(w, r)
//		println("timeout end")
//	})
//}
//
//func main() {
//	http.Handle("/", timeMiddleware(timeOut(http.HandlerFunc(hello))))
//	http.ListenAndServe(":8080", nil)
//}
