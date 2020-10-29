package main

func gameOfLife(board [][]int) {
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[0]); j++ {
			nums := 0
			// 上面
			if i-1 >= 0 {
				nums += board[i-1][j]
			}
			// 左面
			if j-1 >= 0 {
				nums += board[i][j-1]
			}
			// 下面
			if i+1 < len(board) {
				nums += board[i+1][j]
			}
			// 右面
			if j+1 < len(board[i]) {
				nums += board[i][j+1]
			}
			// 左上
			if i-1 >= 0 && j-1 >= 0 {
				nums += board[i-1][j-1]
			}
			// 右上
			if i-1 >= 0 && j+1 < len(board[i]) {
				nums += board[i-1][j+1]
			}
			// 左下
			if i+1 < len(board) && j-1 >= 0 {
				nums += board[i+1][j-1]
			}
			// 右下
			if j+1 < len(board[i]) && i+1 < len(board) {
				nums += board[i+1][j+1]
			}
			switch {
			case nums < 2:
				board[i][j] = 0
			case nums == 3 && board[i][j] == 0:
				board[i][j] = 1
			case nums > 3:
				board[i][j] = 0
			}
			board[i][j] = board[i][j] % 2
		}
	}
}
