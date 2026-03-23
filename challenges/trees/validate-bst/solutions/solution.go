package solution

import "math"

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func validate(node *TreeNode, lo, hi float64) bool {
	if node == nil {
		return true
	}
	val := float64(node.Val)
	if val <= lo || val >= hi {
		return false
	}
	return validate(node.Left, lo, val) && validate(node.Right, val, hi)
}

func IsValidBST(root *TreeNode) bool {
	return validate(root, math.Inf(-1), math.Inf(1))
}
