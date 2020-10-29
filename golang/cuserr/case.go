package cuserr

import (
	"context"
	"net/http"
)

type HandlerMidlleware interface {
	HandleHTTPC(ctx context.Context, rw http.ResponseWriter, req *http.Request, next http.Handler)
}

var function1 HandlerMidlleware
var function2 HandlerMidlleware

func addUserID(rw http.ResponseWriter, req *http.Request, next http.Handler) {
	ctx := context.WithValue(req.Context(), "userid", req.Header.Get("userid"))
	req = req.WithContext(ctx)
	next.ServeHTTP(rw, req)
}
