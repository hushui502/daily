package factory_method

type IRuleConfigParser interface {
	Parse(data []byte)
}

type JsonRuleConfigParser struct {

}

func (j JsonRuleConfigParser) Parse(data []byte) {
	panic("implement me")
}

type YamlRuleConfigParser struct {

}

func (y YamlRuleConfigParser) Parse(data []byte) {
	panic("implement me")
}

type IRuleConfigParserFactory interface {
	CreateParser() IRuleConfigParser
}

type yamlRuleConfigParserFactory struct {
	
}

func (y yamlRuleConfigParserFactory) CreateParser() IRuleConfigParser {
	return YamlRuleConfigParser{}
}

type jsonRuleConfigParserFactory struct {

}

func (j jsonRuleConfigParserFactory) CreateParser() IRuleConfigParser {
	return JsonRuleConfigParser{}
}

func NewIRuleConfigParserFactory(t string) IRuleConfigParserFactory {
	switch t {
	case "json":
		return jsonRuleConfigParserFactory{}
	case "yaml":
		return yamlRuleConfigParserFactory{}
	}

	return nil
}






