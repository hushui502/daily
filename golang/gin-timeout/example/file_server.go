package main

import (
	timeout "gin-timeout"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func main() {
	router := gin.Default()
	defaultMsg := `{"code": -1, "msg":"http: Handler timeout"}`
	router.Use(timeout.Timeout(timeout.WithTimeout(10*time.Second),
		timeout.WithDefaultMsg(defaultMsg)))
	router.Static("/static", "tmp/static")
	log.Fatal(router.Run(":8080"))
}
