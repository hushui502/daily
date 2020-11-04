package main

import "math"

func maxProfit(prices []int) int {
	var minprice = math.MaxInt64
	var maxprofit = 0

	for i := 0; i < len(prices); i++ {
		if prices[i] < minprice {
			minprice = prices[i]
		} else if prices[i]-minprice > maxprofit {
			maxprofit = prices[i] - minprice
		}
	}
	return maxprofit
}
