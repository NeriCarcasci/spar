package solution

import "math"

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func MaxPathSum(root *TreeNode) int {
	result := math.MinInt64
	var dfs func(node *TreeNode) int
	dfs = func(node *TreeNode) int {
		if node == nil {
			return 0
		}
		leftGain := dfs(node.Left)
		if leftGain < 0 {
			leftGain = 0
		}
		rightGain := dfs(node.Right)
		if rightGain < 0 {
			rightGain = 0
		}
		pathSum := node.Val + leftGain + rightGain
		if pathSum > result {
			result = pathSum
		}
		if leftGain > rightGain {
			return node.Val + leftGain
		}
		return node.Val + rightGain
	}
	dfs(root)
	return result
}
