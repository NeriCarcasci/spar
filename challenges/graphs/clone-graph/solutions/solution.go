package solution

type Node struct {
	Val       int
	Neighbors []*Node
}

func CloneGraph(node *Node) *Node {
	if node == nil {
		return nil
	}
	cloned := make(map[*Node]*Node)
	return dfs(node, cloned)
}

func dfs(node *Node, cloned map[*Node]*Node) *Node {
	if c, ok := cloned[node]; ok {
		return c
	}
	copy := &Node{Val: node.Val}
	cloned[node] = copy
	for _, neighbor := range node.Neighbors {
		copy.Neighbors = append(copy.Neighbors, dfs(neighbor, cloned))
	}
	return copy
}
