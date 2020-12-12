package chain

type SensitiveWordFilter interface {
	Filter(content string) bool
}

type SensitiveWordFilterChain struct {
	filters []SensitiveWordFilter
}

func (c *SensitiveWordFilterChain) AddFilter(filter SensitiveWordFilter) {
	c.filters = append(c.filters, filter)
}

func (c *SensitiveWordFilterChain) Filter(content string) bool {
	for _, filter := range c.filters {
		// fast path
		if filter.Filter(content) {
			return true
		}
	}

	return false
}

type AdSensitiveFilter struct {

}

func (a AdSensitiveFilter) Filter(content string) bool {
	return false
}

type PoliticalWordFilter struct {

}

func (p PoliticalWordFilter) Filter(content string) bool {
	return true
}
