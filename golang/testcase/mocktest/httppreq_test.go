package mocktest

import (
	"fmt"
	"testing"
)

func Test_Average(t *testing.T) {
	testError := fmt.Errorf("an example error to compare against")

	tests := []struct {
		name           string
		input          []float64
		err            error
		expectedResult float64
		expectedErr    error
	}{
		{
			name:           "three normal values",
			input:          []float64{3.5, 2.5, 9.0},
			expectedResult: 5,
		},
		{
			name:        "error case",
			input:       []float64{3.5, 2.5, 9.0},
			err:         testError,
			expectedErr: testError,
		},
	}

	for _, test := range tests {
		service := service{valueGetter: mockValueGetter{
			values: test.input,
			err:    test.err,
		}}

		result, err := service.averageForWeb()

		if err != test.expectedErr {
			t.Errorf(
				"for average test %s, got err %v but expected %f",
				test.name,
				test.expectedResult,
				test.expectedErr,
			)
		}

		if result != test.expectedResult {
			//...
			// 应该在这里设置expecteErr
			t.Errorf(
				"for average test %s, got result %f but expected %f",
				test.name,
				result,
				test.expectedResult,
			)
		}
	}
}
