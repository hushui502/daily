package leetcode

func maxProfit3(prices []int) int {
	profit := 0

	for i := 0; i < len(prices); i++ {
		if prices[i+1] > prices[i] {
			profit += prices[i+1] - prices[i]
		}
	}

	return profit
}
