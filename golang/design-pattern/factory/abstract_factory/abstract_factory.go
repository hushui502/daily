package abstract_factory

type IRuleConfigParser interface {
	Parse(data []byte)
}

type ISystemConfigParser interface {
	ParseSystem(data []byte)
}

type IConfigParserFactory interface {
	CreateRuleParser() IRuleConfigParser
	CreateSystemParser() ISystemConfigParser
}

type jsonRuleConfigParser struct {}

func (j jsonRuleConfigParser) Parse(data []byte) {
	panic("implement me")
}

type jsonSystemConfigParser struct {}

func (j jsonSystemConfigParser) ParseSystem(data []byte) {
	panic("implement me")
}

type jsonConfigParserFactory struct {}

func (j jsonConfigParserFactory) CreateRuleParser() IRuleConfigParser {
	return jsonRuleConfigParser{}
}

func (j jsonConfigParserFactory) CreateSystemParser() ISystemConfigParser {
	return jsonSystemConfigParser{}
}



