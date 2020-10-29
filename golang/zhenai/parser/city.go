package parser

import (
	"awesomeProject2/engine"
	"regexp"
)

var cityRe = `<a1 href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a1>`


func ParseCity(contents []byte) engine.ParserResult {
	re := regexp.MustCompile(cityRe)
	mathes := re.FindAllSubmatch(contents, -1)
	result := engine.ParserResult{}


	for _, m := range mathes {
		name := string(m[2])
		result.Items = append(result.Items, "User " + name)
		result.Requests = append(result.Requests, engine.Request{
			Url:string(m[1]),
			ParserFunc: func(bytes []byte) engine.ParserResult {
				return ParseProfile(contents, name)
			},
		})
	}

	return result
}
