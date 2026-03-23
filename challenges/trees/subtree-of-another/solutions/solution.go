package solution

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func isSame(p *TreeNode, q *TreeNode) bool {
	if p == nil && q == nil {
		return true
	}
	if p == nil || q == nil {
		return false
	}
	return p.Val == q.Val && isSame(p.Left, q.Left) && isSame(p.Right, q.Right)
}

func IsSubtree(root *TreeNode, subRoot *TreeNode) bool {
	if root == nil {
		return false
	}
	if isSame(root, subRoot) {
		return true
	}
	return IsSubtree(root.Left, subRoot) || IsSubtree(root.Right, subRoot)
}
