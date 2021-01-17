package leetcode

// dp[i] = dp[i-1] + nums[i]
// if dp[i-1] < 0 ? dp[i] = nums[i]. So dp[i] = max(dp[i-1]+nums[i], nums[i])

func maxSubArray(nums []int) int {
	if len(nums) < 1 {
		return 0
	}

	dp := make([]int, len(nums))
	result := nums[0]
	dp[0] = nums[0]

	for i := 1; i < len(nums); i++ {
		dp[i] = max(dp[i-1]+nums[i], nums[i])
		result = max(result, dp[i])
	}

	return result
}
