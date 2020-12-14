package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/testdata/protoexample"
	"net/http"
)

type Person struct {
	ID string `uri:"id" binding:"required,uuid"`
	Name string `uri:"name" binding:"required"`
}

type myForm struct {
	Colors []string `form:"colors[]"`
}


func main() {
	route := gin.Default()
	//route.GET("/:name/:id", func(c *gin.Context) {
	//	var person Person
	//	if err := c.ShouldBindUri(&person); err != nil {
	//		c.JSON(400, gin.H{"msg": err})
	//		return
	//	}
	//	c.JSON(200, gin.H{"name": person.Name, "uuid": person.ID})
	//})
	route.GET("/someProtoBuf", func(c *gin.Context) {
		resp := []int64{int64(1), int64(2)}
		label := "test"
		data := &protoexample.Test{
			Label:&label,
			Reps:resp,
		}
		c.ProtoBuf(200, data)
	})
	route.GET("/someJSON", func(c *gin.Context) {
		names := []string{"lena", "austin", "foo"}

		// 将会输出:   while(1);["lena","austin","foo"]
		c.SecureJSON(http.StatusOK, names)
	})
	route.GET("/json", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"html": "<b>Hello, world!</b>",
		})
	})

	route.GET("/purejson", func(c *gin.Context) {
		c.PureJSON(200, gin.H{
			"html": "<b>Hello, world!</b>",
		})
	})
	route.Run(":8088")
}