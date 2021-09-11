package main

type IRuleConfigParser interface {
	Parse(data []byte)
}

type jsonParser struct {}

func (j jsonParser) Parse(data []byte) {
	panic("success json")
}

type yamlParser struct {}

func (y yamlParser) Parse(data []byte) {
	panic("success yaml")
}

func NewParser(t string) IRuleConfigParser {
	switch t {
	case "json":
		return jsonParser{}
	case "yaml":
		return yamlParser{}
	}
	return nil
}

