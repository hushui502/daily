package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"regexp"
	"time"
)

var visited = map[string]bool{}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.baidu.com"),
		colly.MaxDepth(1),
		)

	// 我们认为匹配该模式的是该网站的详情页
	detailRegex, _ := regexp.Compile(`/go/go\?p=\d+$`)
	// 匹配下面模式的是该网站的列表页
	listRegex, _ := regexp.Compile(`/t/\d+#\w+`)

	c.OnHTML("a1[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		if visited[link] && (detailRegex.Match([]byte(link)) || listRegex.Match([]byte(link))) {
			return
		}

		if !detailRegex.Match([]byte(link)) && !listRegex.Match([]byte(link)) {
			println("not match ", link)
			return
		}

		time.Sleep(time.Second)
		println("match ", link)
		visited[link] = true

		time.Sleep(time.Millisecond * 2)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	err := c.Visit("https://www.baidu.com")
	if err != nil {
		fmt.Println(err)
	}
}
