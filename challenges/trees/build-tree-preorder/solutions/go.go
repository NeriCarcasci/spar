package solution

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func BuildTree(preorder []int, inorder []int) *TreeNode {
	if len(preorder) == 0 {
		return nil
	}
	rootVal := preorder[0]
	mid := 0
	for i, v := range inorder {
		if v == rootVal {
			mid = i
			break
		}
	}
	return &TreeNode{
		Val:   rootVal,
		Left:  BuildTree(preorder[1:mid+1], inorder[:mid]),
		Right: BuildTree(preorder[mid+1:], inorder[mid+1:]),
	}
}
