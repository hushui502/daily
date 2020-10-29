package main

import (
	"awesomeProject2/gee-web/group/jii"
	"fmt"
	"net/http"
	"time"
)

type student struct {
	Name string
	Age  int
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := jii.Default()
	r.GET("/", func(c *jii.Context) {
		c.String(http.StatusOK, "Hello Jii\n")
	})
	r.GET("/panic", func(c *jii.Context) {
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[100])
	})
	r.Run(":9999")
	//r := jii.New()
	//r.Use(jii.Logger())
	//r.SetFuncMap(template.FuncMap{
	//	"formatAsDate":formatAsDate,
	//})
	//r.LoadHTMLGlob("D:\\project\\go\\src\\awesomeProject2\\jii-web\\group\\templates\\*")
	//r.Static("/assets", "./static")
	//
	//stu1 := &student{Name:"hufan", Age:18}
	//stu2 := &student{Name:"wangzhen", Age:17}
	//r.GET("/", func(c *jii.Context) {
	//	c.HTML(http.StatusOK, "css.tmpl", nil)
	//})
	//r.GET("/students", func(c *jii.Context) {
	//	c.HTML(http.StatusOK, "arr.tmpl", jii.H{
	//		"title": "jii",
	//		"stuArr": [2]*student{stu1, stu2},
	//	})
	//})
	//r.GET("/date", func(c *jii.Context) {
	//	c.HTML(http.StatusOK, "custom_func.tmpl", jii.H{
	//		"title": "jii",
	//		"now" : time.Date(2019, 3, 3, 0,0,0,0,time.UTC),
	//	})
	//})
	//r.Run(":9999")
}
