package solution

import (
	"strconv"
	"strings"
)

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func Serialize(root *TreeNode) string {
	var parts []string
	var dfs func(node *TreeNode)
	dfs = func(node *TreeNode) {
		if node == nil {
			parts = append(parts, "N")
			return
		}
		parts = append(parts, strconv.Itoa(node.Val))
		dfs(node.Left)
		dfs(node.Right)
	}
	dfs(root)
	return strings.Join(parts, ",")
}

func Deserialize(data string) *TreeNode {
	vals := strings.Split(data, ",")
	idx := 0
	var dfs func() *TreeNode
	dfs = func() *TreeNode {
		if vals[idx] == "N" {
			idx++
			return nil
		}
		val, _ := strconv.Atoi(vals[idx])
		idx++
		node := &TreeNode{Val: val}
		node.Left = dfs()
		node.Right = dfs()
		return node
	}
	return dfs()
}
