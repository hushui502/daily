package leetcode

func findPathSum(grid [][]int) int {
	l := len(grid)
	if l < 1 {
		return 0
	}
	dp := make([][]int, l)
	for i, arr := range grid {
		dp[i] = make([]int, len(arr))
	}
	dp[0][0] = grid[0][0]
	for i := 0; i < l; i++ {
		for j := 0; j < len(grid[i]); j++ {
			if i == 0 && j != 0 {
				dp[i][j] = dp[i][j-1] + grid[i][j]
			} else if i != 0 && j == 0 {
				dp[i][j] = dp[i-1][j] + grid[i][j]
			} else if i != 0 && j != 0 {
				dp[i][j] = min(dp[i-1][j], dp[i][j-1]) + grid[i][j]
			}
		}
	}

	return dp[l-1][len(dp[l-1])-1]
}
