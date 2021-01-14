package leetcode


// 状态转移方程
func maxProfit(princes []int) int {
	n := len(princes)
	dp := make([][2]int, n)
	dp[0][1] = -princes[0]
	for i := 1; i < n; i++ {
		// 如果当前天没有股票,说明上一天就没有,或者上一天有今天卖掉(收益加prices[i])
		dp[i][0] = max(dp[i-1][0], dp[i-1][1]+princes[i])
		// 如果当前天有股票,说明上一天就有,或者上一天没有今天买入
		dp[i][1] = max(dp[i-1][1], dp[i-1][0]-princes[i])
	}

	return dp[n-1][0]
}
