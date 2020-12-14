package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

func main() {
	router := gin.Default()
	router.GET("/some", func(c *gin.Context) {
		response, err := http.Get("https://raw.githubusercontent.com/gin-gonic/logo/master/color.png")
		if err != nil || response.StatusCode != http.StatusOK {
			c.Status(http.StatusServiceUnavailable)
			return
		}
		reader := response.Body
		contentLength := response.ContentLength
		contentType := response.Header.Get("Content-Type")
		extraHeaders := map[string]string{
			"Content-Disposition": `attachment; filename="gopher.png"`,
		}
		c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
	})

	router.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})

	router.LoadHTMLFiles("./")
	router.LoadHTMLGlob("")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "", )
	})

	router.Run(":8080")

}

func formatAsDate() string {
	return fmt.Sprintf("")
}
