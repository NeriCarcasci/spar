package solution

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func LowestCommonAncestor(root *TreeNode, p int, q int) int {
	cur := root
	for cur != nil {
		if p < cur.Val && q < cur.Val {
			cur = cur.Left
		} else if p > cur.Val && q > cur.Val {
			cur = cur.Right
		} else {
			return cur.Val
		}
	}
	return -1
}
