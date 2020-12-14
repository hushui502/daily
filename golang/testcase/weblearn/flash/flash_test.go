package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCookie(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/post", GetCookie)

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "post/1", nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {

	}
	var post Post
	json.Unmarshal(writer.Body.Bytes(), &post)

}
