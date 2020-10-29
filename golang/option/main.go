package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)
type Profile struct {
	Status    string
	Message	 string
	Data string
}
type TT struct {
	FilePath string `json:"file_path"`
}
func main() {
	http.HandleFunc("/", foo)
	http.ListenAndServe(":3000", nil)
}

func foo(w http.ResponseWriter, r *http.Request) {
	//file_path := "sss"
	//tt := TT{FilePath:`{"encryption":"md5","timestamp":1482463793,"key":"2342874840784a81d4d9e335aaf76260","partnercode":100034}`}

	tt := TT{FilePath:"/tmp/path"}
	tts, _ := json.Marshal(tt)
	fmt.Println(string(tts))
	//tmpStr := string(tts)
	//strings.ReplaceAll(tmpStr, ""\", "")
	//res := json.RawMessage(string(tts))

	profile := Profile{
		Status:"success",
		Message:"upload file successed",
		Data: string(tts),
	}
	js, err := json.Marshal(profile)
	js = bytes.ReplaceAll(js, []byte(`\`), []byte(``))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}