package parser

import (
	"awesomeProject2/engine"
	"awesomeProject2/model"
	"regexp"
	"strconv"
)

var (
	UrlIdMatcher         = regexp.MustCompile(`http://album.zhenai.com/u/([\d]+)`)
	AgeMatcher           = regexp.MustCompile(`<td><span[^>]*>年龄：</span>([\d]+)岁</td>`)
	GenderMatcher        = regexp.MustCompile(`<td><span[^>]*>性别：</span><span field="">([^<]+)</span></td>`)
	MarriageMatcher      = regexp.MustCompile(`<td><span[^>]*>婚况：</span>([^<]+)</td>`)
	HeightMatcher        = regexp.MustCompile(`<td><span[^>]*>身高：</span><span field="">([\d]+)CM</span></td>`)
	WeightMatcher        = regexp.MustCompile(`<td><span[^>]*>体重：</span><span field="">([\d]+)</span></td>`)
	IncomeMatcher        = regexp.MustCompile(`<td><span[^>]*>月收入：</span>([^<]+)</td>`)
	EducationMatcher     = regexp.MustCompile(`<td><span[^>]*>学历：</span>([^<]+)</td>`)
	OccupationMatcher    = regexp.MustCompile(`<td><span[^>]*>职业： </span>([^<]+)</td>`)
	ConstellationMatcher = regexp.MustCompile(`<td><span[^>]*>星座：</span><span field="">([^<]+)</span></td>`)
	HouseMatcher         = regexp.MustCompile(`<td><span[^>]*>住房条件：</span><span field="">([^<]+)</span></td>`)
	CarMatcher           = regexp.MustCompile(`<td><span[^>]*>是否购车：</span><span field="">([^<]+)</span></td>`)
	RecommendMatcher     = regexp.MustCompile(`<a1 class="exp-user-name"[^>]*href="(http://album.zhenai.com/u/[\d]+)">([^<]+)</a1>`)
)

func ParseProfile(contents []byte, name string) engine.ParserResult {
	profile := model.Profile{}
	age, err := strconv.Atoi(extractString(contents, AgeMatcher))

	profile.Name = name
	if err == nil {
		profile.Age = age
	}
	profile.Marriage = extractString(contents, MarriageMatcher)
	profile.Occupation = extractString(contents, OccupationMatcher)
	profile.House = extractString(contents, HouseMatcher)
	profile.Car = extractString(contents, CarMatcher)
	profile.Gender = extractString(contents, GenderMatcher)
	profile.Income = extractString(contents, IncomeMatcher)
	profile.Education = extractString(contents, EducationMatcher)

	result := engine.ParserResult{
		Items: []interface{}{profile},
	}

	return result
}

func extractString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)
	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}