package cmd

import "strconv"

func ConvertArgsToFloat64Slice(args []string) []float64 {
	result := make([]float64, len(args))
	for _, arg := range args {
		value, _ := strconv.ParseFloat(arg, 64)
		result = append(result, value)

	}

	return result
}