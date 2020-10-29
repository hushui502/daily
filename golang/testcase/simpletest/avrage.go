package simpletest

func average(values []float64) float64 {
	var total float64 = 0
	for _, value := range values {
		total += value
	}
	return total / float64(len(values))
}
