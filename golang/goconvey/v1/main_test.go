package main

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestHelloHandleFunc(t *testing.T) {
	rw := httptest.NewRecorder()
	name := "hufan"
	req := httptest.NewRequest(http.MethodPost, "/hello?name="+name, nil)
	handleHello(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("status code not ok, status code is %v", rw.Code)
	}

	if len(counter) != 1 {
		t.Errorf("counter len not correct")
	}

	if counter[name] != 1 {
		t.Errorf("counter value is error: visitor=%s count=%v", name, counter[name])
	}
}

func TestHTTPServer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handleHello))
	defer ts.Close()

	logrus.Infof("server url: %s", ts.URL)

	testURL := ts.URL + "/hello?name=hufan"
	resp, err := http.Get(testURL)
	if err != nil {
		t.Error(err)
		return
	}
	if g, w := resp.StatusCode, http.StatusOK; g != w {
		t.Errorf("status code = %q; want %q", g, w)
		return
	}
}

func TestHelloHandlerInSubtest(t *testing.T) {

	tests := []struct {
		name string
		wCnt int
	}{
		{name: "zouying", wCnt: 1},
		{name: "user2", wCnt: 1},
		{name: "user3", wCnt: 1},
	}

	for _, tc := range tests {
		t.Run("test-"+tc.name, func(t *testing.T) {
			rw := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/hello?name="+tc.name, nil)
			handleHello(rw, req)

			assert.Equal(t, http.StatusOK, rw.Code)
			assert.Equal(t, tc.wCnt, counter[tc.name])
		})
	}
}

func TestHelloHandlerDetectDataRace(t *testing.T) {
	tests := []struct {
		name string
		wCnt int
	}{
		{name: "zouying", wCnt: 1},
		{name: "zouying", wCnt: 2},
		{name: "user2", wCnt: 1},
		{name: "user3", wCnt: 1},
	}

	var wg sync.WaitGroup
	wg.Add(len(tests))
	for _, tc := range tests {
		name := tc.name
		want := tc.wCnt

		go func() {
			defer wg.Done()

			rw := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/hello?name="+name, nil)
			handleHello(rw, req)

			assert.Equal(t, http.StatusOK, rw.Code)
			assert.Equal(t, want, counter[name])
		}()
	}
	wg.Wait()
}
