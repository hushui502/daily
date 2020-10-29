package mocktest

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ValueGetter interface {
	GetValues() ([]float64, error)
}

type service struct {
	valueGetter ValueGetter
}

func (s service) averageForWeb() (float64, error) {
	values, err := s.valueGetter.GetValues()
	if err != nil {
		return 0, err
	}

	var total float64 = 0
	for _, value := range values {
		total += value
	}
	return total / float64(len(values)), nil
}

// The real implementation

type httpValueGetter struct {
}

func (h httpValueGetter) GetValues() ([]float64, error) {
	resp, err := http.DefaultClient.Get("http://our-api.com/numbers")
	if err != nil {
		return nil, err
	}

	values := []float64{}
	if err := json.NewDecoder(resp.Body).Decode(&values); err != nil {
		return nil, nil
	}
	return values, nil
}

type mockValueGetter struct {
	values []float64
	err    error
}

func (m mockValueGetter) GetValues() ([]float64, error) {
	return m.values, m.err
}

func main() {
	service := service{valueGetter: httpValueGetter{}}

	average, err := service.averageForWeb()
	if err != nil {
		panic(err)
	}

	fmt.Println(average)
}
