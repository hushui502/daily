package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

func handleSearch(w http.ResponseWriter, r *http.Request) {
	var (
		ctx context.Context
		cancel context.CancelFunc
	)

	timeout, err := time.ParseDuration(r.FormValue("timeout"))
	if err != nil {
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	}
	defer cancel()

	// check the search query
	query := r.FormValue("q")
	if err != nil {
		http.Error(w, "no query", http.StatusBadRequest)
		return
	}

	// store the user ip in ctx for user by code in other packages.
	userIp, err := FromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx = NewContext(ctx, userIp)

	// run the search and print the results.
	//start := time.Now()
	results, _ := Search(ctx, query)
	//elapsed := time.Since(start)

	res, _ := json.Marshal(results)
	w.Write(res)

}

/* ============userip==================*/
func FromRequest(req *http.Request) (net.IP, error) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return nil, err
	}

	userIp := net.ParseIP(ip)
	if userIp == nil {
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}

	return userIp, nil
}

// The key is unexported to protect collisions with context keys defind in
// other packages.
type key int

const userIPKey = 0

func NewContext(ctx context.Context, userIp net.IP) context.Context {
	return context.WithValue(ctx, userIPKey, userIp)
}

func FromContext(ctx context.Context) (net.IP, bool) {
	userIP, ok := ctx.Value(userIPKey).(net.IP)

	return userIP, ok
}

/* =====================search===================== */
type Results []Result

type Result struct {
	Title, URL string
}

func Search(ctx context.Context, query string) (Results, error) {
	// Prepare the Goggle Search API request.
	req, err := http.NewRequest("GET", "https://ajax.googleapis.com/ajax/services/search/web?v=1.0", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("q", query)

	if userIP, ok := FromContext(ctx); ok {
		q.Set("userip", userIP.String())
	}
	req.URL.RawQuery = q.Encode()

	var results Results
	err = httpDo(ctx, req, func(resp *http.Response, err error) error {
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// parse the json search results
		var data struct{
			ResponseData struct{
				Results []struct{
					TitleNoFormatting string
					URL string
				}
			}
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return err
		}
		for _, res := range data.ResponseData.Results {
			results = append(results, Result{Title: res.TitleNoFormatting, URL: res.URL})
		}

		return nil
	})

	// httpDo waits for the closure we provided to returns, so it is safe to
	// read result here.

	return results, err
}

func httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	c := make(chan error, 1)
	req = req.WithContext(ctx)
	go func() {
		c <- f(http.DefaultClient.Do(req))
	}()
	
	select {
	case <-ctx.Done():
		<-c		// wait for to return
		return ctx.Err()
	case err := <-c:
		return err
	}
}