package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckStatusOK(t *testing.T) {
	req, err := http.NewRequest("GET", "/CheckStatusOK", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CheckStatusOK)
	handler.ServeHTTP(rr, req)
	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("handler returned %v", status)
	}

	expect := `fine!`
	if rr.Body.String() != expect {
		t.Errorf("handler returned %v", rr.Body.String())
	}
}
