// https://github.com/vearne/gin-timeout
package gin_timeout

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var (
	defaultOptions TimeoutOptions
)

func init() {
	defaultOptions = TimeoutOptions{
		CallBack:      nil,
		DefaultMsg:    `{"code": -1, "msg":"http: Handler timeout"}`,
		Timeout:       time.Duration(3) * time.Second,
		ErrorHttpCode: http.StatusServiceUnavailable,
	}
}

func Timeout(opts ...Option) gin.HandlerFunc {
	return func(c *gin.Context) {
		buffer := GetBuff()

		tw := &TimeoutWriter{
			body:           buffer,
			ResponseWriter: c.Writer,
			h:              make(http.Header),
		}

		tw.TimeoutOptions = defaultOptions

		// loop through each option
		for _, opt := range opts {
			opt(tw)
		}
		c.Writer = tw

		// wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), tw.Timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)
		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			c.Next()
			finish <- struct{}{}
		}()

		var err error
		select {
		case p := <-panicChan:
			panic(p)
		case <-ctx.Done():
			tw.mu.Lock()
			defer tw.mu.Unlock()

			tw.timedOut = true
			tw.ResponseWriter.WriteHeader(tw.ErrorHttpCode)
			_, err = tw.ResponseWriter.Write([]byte(tw.DefaultMsg))
			if err != nil {
				panic(err)
			}
			c.Abort()

			if tw.CallBack != nil {
				tw.CallBack(c.Request.Clone(context.Background()))
			}
		case <-finish:
			tw.mu.Lock()
			defer tw.mu.Unlock()
			dst := tw.ResponseWriter.Header()
			for k, v := range tw.Header() {
				dst[k] = v
			}

			if !tw.wroteHeader {
				tw.code = http.StatusOK
			}

			tw.ResponseWriter.WriteHeader(tw.code)
			_, err = tw.ResponseWriter.Write(buffer.Bytes())
			if err != nil {
				panic(err)
			}
			PutBuff(buffer)
		}
	}
}
