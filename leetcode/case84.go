package leetcode

func sortArrayToBSF(nums []int) *TreeNode {
	if len(nums) == 0 {
		return nil
	}

	return &TreeNode{Val:nums[len(nums)/2], Left:sortArrayToBSF(nums[:len(nums)/2]), Right:sortArrayToBSF(nums[len(nums)/2+1:])}
}
