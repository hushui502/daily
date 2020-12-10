package leetcode

func maxProfit(prices []int) int {
	if len(prices) < 1 {
		return 0
	}
	min, maxProfit := prices[0], 0
	
	for i := 1; i < len(prices); i++ {
		if prices[i] - min > maxProfit {
			maxProfit = prices[i] - min
		}
		if prices[i] < min {
			min = prices[i]
		}
	}

	return maxProfit
}

func maxProfit2(prices []int) int {
	if len(prices) < 1 {
		return 0
	}

	stack, res := []int{prices[0]}, 0
	for i := 0; i < len(prices); i++ {
		if prices[i] < stack[len(stack)-1] {
			index := len(stack) - 1
			for ; index >= 0; index-- {
				if stack[index] < prices[i] {
					break
				}
			}
			stack = stack[:index+1]
		}

		stack = append(stack, prices[i])
		res = max(res, stack[len(stack)-1] - stack[0])
	}
}
