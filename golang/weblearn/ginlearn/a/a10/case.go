package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func main() {
	r := gin.New()

	r.GET("/login_async", func(c *gin.Context) {
		ccp := c.Copy()
		go func() {
			time.Sleep(2 * time.Second)
			log.Println("done! " + ccp.Request.URL.Path)
		}()
	})

	r.GET("/long_sync", func(c *gin.Context) {
		// simulate a long task with time.Sleep(). 5 seconds
		time.Sleep(5 * time.Second)

		// 这里没有使用goroutine，所以不用使用副本
		log.Println("Done! in path " + c.Request.URL.Path)
	})

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}
