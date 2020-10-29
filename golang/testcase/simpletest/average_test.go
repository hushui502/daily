package simpletest

import "testing"

func Test_Average(t *testing.T) {
	tests := []struct {
		name           string
		input          []float64
		expectedResult float64
	}{
		{
			name:           "case 1",
			input:          []float64{3.5, 2.5, 9.0},
			expectedResult: 5,
		},
		{
			name:           "case 2",
			input:          []float64{0, 0},
			expectedResult: 0,
		},
	}

	for _, test := range tests {
		result := average(test.input)
		if result != test.expectedResult {
			t.Errorf(
				"for average test %s, got result %f but expected %f",
				test.name,
				result,
				test.expectedResult,
			)
		}
	}
}
