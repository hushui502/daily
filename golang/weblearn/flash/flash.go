package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}

func Log(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
		fmt.Println("Log---", name)
		h(w, r)
	}
}

func Header(w http.ResponseWriter, r *http.Request) {
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	fmt.Fprintln(w, string(body))
}

func Process(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024)
	fileHeader := r.MultipartForm.File["uploaded"][0]
	file, err := fileHeader.Open()
	if err == nil {
		data, err := ioutil.ReadAll(file)
		if err == nil {
			fmt.Fprintln(w, string(data))
		}
	}
}

func HeaderEx(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "https://www.baidu.com")
	w.WriteHeader(302)
}

type Post struct {
	User    string
	Threads []string
}

func JsonEx(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	post := &Post{
		User:    "hufan",
		Threads: []string{"1,", "4r4", "6565"},
	}
	res, _ := json.Marshal(post)
	w.Write(res)
}

func SetCookie(w http.ResponseWriter, r *http.Request) {
	c1 := http.Cookie{Name: "first", Value: "go web", HttpOnly: true}
	c2 := http.Cookie{Name: "secend", Value: "java web", HttpOnly: true}

	w.Header().Set("Set-Cookie", c1.String())
	w.Header().Add("Set-Cookie", c2.String())
	http.SetCookie(w, &c1)
	str, _ := strconv.Atoi(path.Base(r.URL.Path))
	fmt.Println(str)
}

func GetCookie(w http.ResponseWriter, r *http.Request) {
	c1, _ := r.Cookie("first")
	cs := r.Cookies()
	fmt.Fprintln(w, c1)
	fmt.Fprintln(w, cs)
}

func SetMessage(w http.ResponseWriter, r *http.Request) {
	msg := []byte("Hello")
	c := http.Cookie{
		Name:  "flash",
		Value: base64.URLEncoding.EncodeToString(msg),
	}
	http.SetCookie(w, &c)
}

func GetMessage(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("flash")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Fprintln(w, "No message cookie")
		}
	} else {
		rc := http.Cookie{
			Name:    "flash",
			MaxAge:  -1,
			Expires: time.Unix(1, 0),
		}
		http.SetCookie(w, &rc)
		val, _ := base64.URLEncoding.DecodeString(c.Value)
		fmt.Fprintln(w, string(val))
	}
}

func main() {
	server := http.Server{Addr: "127.0.0.1:8080"}
	http.HandleFunc("/Hello", SetCookie)
	http.HandleFunc("/get", GetCookie)
	http.HandleFunc("/SetMessage", SetMessage)
	http.HandleFunc("/GetMessage", GetMessage)
	//http2.ConfigureServer(&server, &http2.Server{})

	server.ListenAndServe()
}
