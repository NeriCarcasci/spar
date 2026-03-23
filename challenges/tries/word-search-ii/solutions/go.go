package solution

import "sort"

type trieNode struct {
	children [26]*trieNode
	word     string
}

func FindWords(board [][]byte, words []string) []string {
	root := &trieNode{}
	for _, w := range words {
		node := root
		for _, ch := range w {
			idx := ch - 'a'
			if node.children[idx] == nil {
				node.children[idx] = &trieNode{}
			}
			node = node.children[idx]
		}
		node.word = w
	}

	rows, cols := len(board), len(board[0])
	var result []string
	dirs := [4][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

	var dfs func(r, c int, node *trieNode)
	dfs = func(r, c int, node *trieNode) {
		ch := board[r][c]
		if ch == '#' {
			return
		}
		next := node.children[ch-'a']
		if next == nil {
			return
		}
		if next.word != "" {
			result = append(result, next.word)
			next.word = ""
		}
		board[r][c] = '#'
		for _, d := range dirs {
			nr, nc := r+d[0], c+d[1]
			if nr >= 0 && nr < rows && nc >= 0 && nc < cols {
				dfs(nr, nc, next)
			}
		}
		board[r][c] = ch
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			dfs(r, c, root)
		}
	}
	sort.Strings(result)
	return result
}
