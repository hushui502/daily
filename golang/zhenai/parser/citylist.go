package parser

import (
	"awesomeProject2/engine"
	"regexp"
)

const cityListRe  = `<a1 href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a1>`
func ParseCityList(contents []byte) engine.ParserResult {
	re := regexp.MustCompile(cityListRe)
	mathes := re.FindAllSubmatch(contents, -1)
	result := engine.ParserResult{}

	for _, m := range mathes {
		result.Items = append(result.Items, "City " + string(m[2]))
		result.Requests = append(result.Requests, engine.Request{
			Url:string(m[1]),
			ParserFunc:ParseCity,
		})
	}

	return result
}


