package pattern

type Seller interface {
	sell(name string)
}

type Station struct {
	stock int
}

func (s *Station) sell(nam string) {
	if s.stock > 0 {
		s.stock--
	} else {
		
	}
}

type StationProxy struct {
	station *Station
}

func (sp *StationProxy) sell(name string) {
	if sp.station.stock > 0 {
		sp.station.stock--
	} else {

	}
}

func testproxy() {
	station := &Station{3}
	proxy := StationProxy{station}
	station.sell("a")
	proxy.sell("b")
}
