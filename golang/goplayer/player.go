package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Entry struct {
	Name string // name of the object
	IsDir bool
	Mode os.FileMode
}

const (
	filePrefix = "/f/"
)

var (
	addr = flag.String("http", ":8000", "http listen address")
	root = flag.String("root", "practice/goplayer/", "music root")
)

func main() {
	flag.Parse()
	http.HandleFunc("/", Index)
	http.HandleFunc(filePrefix, File)
	http.ListenAndServe(*addr, nil)
}

func Index(w http.ResponseWriter, r *http.Request) {
	// todo setting
	http.ServeFile(w, r, "practice/goplayer/index.html")
	log.Print("index called")
}

func File(w http.ResponseWriter, r *http.Request) {
	fn := filepath.Join(*root, r.URL.Path[len(filePrefix):])
	fi, err := os.Stat(fn)
	log.Print("File called: ", fn)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if fi.IsDir() {
		serveDirectory(fn, w, r)
		return
	}

	http.ServeFile(w, r, fn)
}

func serveDirectory(fn string, w http.ResponseWriter, r *http.Request) {
	defer func() {
		// 这里要ok不要直接err != nil， 否则interface是一个nil 不能强转成error
		if err, ok := recover().(error); ok {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}()
	d, err := os.Open(fn)
	if err != nil {
		panic(err)
	}
	defer d.Close()

	log.Print("serverDirectory called: ", fn)

	files, err := d.Readdir(-1)
	if err != nil {
		panic(err)
	}
	log.Print("---------", len(files))

	// json encode isn't working with the FileInfo interface,
	// therefore populate an array of Entry and add the Name method
	entries := make([]Entry, len(files), len(files))

	for k := range files {
		log.Print(files[k].Name())
		entries[k].Name = files[k].Name()
		entries[k].IsDir = files[k].IsDir()
		entries[k].Mode = files[k].Mode()
	}

	j := json.NewEncoder(w)

	if err := j.Encode(&entries); err != nil {
		panic(err)
	}
}

