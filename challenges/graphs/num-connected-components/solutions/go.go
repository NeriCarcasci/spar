package solution

func CountComponents(n int, edges [][]int) int {
	parent := make([]int, n)
	rank := make([]int, n)
	for i := range parent {
		parent[i] = i
	}

	find := func(x int) int {
		for parent[x] != x {
			parent[x] = parent[parent[x]]
			x = parent[x]
		}
		return x
	}

	components := n
	for _, e := range edges {
		px, py := find(e[0]), find(e[1])
		if px != py {
			if rank[px] < rank[py] {
				px, py = py, px
			}
			parent[py] = px
			if rank[px] == rank[py] {
				rank[px]++
			}
			components--
		}
	}
	return components
}
