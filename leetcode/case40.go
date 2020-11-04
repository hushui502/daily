package main

func spiralOrder2(matrix [][]int) []int {
	m := len(matrix)
	if m == 0 {
		return nil
	}

	n := len(matrix[0])
	if n == 0 {
		return nil
	}

	top, left, bottom, right := 0, 0, m-1, n-1
	count, sum := 0, m*n
	res := []int{}

	for count < sum {
		i, j := top, left
		for j <= right && count < sum {
			res = append(res, matrix[i][j])
			count++
			j++
		}
		i, j = top+1, right
		for i <= bottom && count < sum {
			res = append(res, matrix[i][j])
			count++
			i++
		}
		i, j = bottom, right-1
		for j >= left && count < sum {
			res = append(res, matrix[i][j])
			count++
			j++
		}
		i, j = bottom-1, left
		for i > top && count < sum {
			res = append(res, matrix[i][j])
			count++
			i--
		}
		top, left, bottom, right = top+1, left+1, bottom-1, right-1
	}

	return res
}