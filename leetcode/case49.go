package leetcode

func minPathSum(grid [][]int) int {
	m, n := len(grid), len(grid[0])
	for i := 0; i < m; i++ {
		grid[i][0] += grid[i-1][0]
	}
	for j := 0; j < n; j++ {
		grid[0][j] += grid[0][j-1]
	}

	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			grid[i][j] += min(grid[i-1][0], grid[0][j-1])
		}
	}

	return grid[m-1][n-1]
}

func minPathSum1(grid [][]int) int {
	if len(grid) == 0 {
		return 0
	}
	m, n := len(grid), len(grid[0])
	if m == 0 || n == 0 {
		return 0
	}

	dp := make([][]int, m)
	for i := 0; i < m; i++ {
		dp[i] = make([]int, n)
	}

	for i := 0; i < len(dp[0]); i++ {
		if i == 0 {
			dp[0][i] = grid[0][i]
		} else {
			dp[0][i] = grid[0][i] + dp[0][i-1]
		}
	}
	for j := 0; j < len(dp); j++ {
		if j == 0 {
			dp[j][0] = dp[j][0]
		} else {
			dp[j][0] = grid[j][0] + grid[j-1][0]
		}
	}

	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			dp[i][j] = min(dp[i-1][j], dp[i][j-1]) + grid[i][j]
		}
	}

	return dp[m-1][n-1]
}

func min(a, b int) int {
	if a >= b {
		return b
	} else {
		return a
	}
}
