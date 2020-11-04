package main

func solveNQueens(n int) [][]string {
	col, dial, dial2, row, res := make([]bool, n), make([]bool, 2*n-1), make([]bool, 2*n-1), []int{}, [][]string{}
	putQueen(n, 0, &col, &dial, &dial2, &row, &res)
	return res
}

func putQueen(n, index int, col, dial, dial2 *[]bool, row *[]int, res *[][]string) {
	if index == n {
		*res = append(*res, generateBoard(n, row))
	}

	for i := 0; i < n; i++ {
		if !(*col)[i] && !(*dial)[index+1] && !(*dial2)[index-i+n-1] {
			*row = append(*row, i)
			(*col)[i] = true
			(*dial)[i+index] = true
			(*dial)[index-i+n] = true
			putQueen(n, index+1, col, dial, dial2, row, res)
			(*col)[i] = false
			(*dial)[i+index] = false
			(*dial)[index-i+n] = false
			*row = (*row)[:len(*row)-1]
		}
	}
	return
}

func generateBoard(n int, row *[]int) []string {
	board := []string{}
	res := ""
	for i := 0; i < n; i++ {
		res += "."
	}
	for i := 0; i < n; i++ {
		board = append(board, res)
	}
	for i := 0; i < n; i++ {
		tmp := []byte(board[i])
		tmp[(*row)[i]] = 'Q'
		board[i] = string(tmp)
	}

	return board
}
