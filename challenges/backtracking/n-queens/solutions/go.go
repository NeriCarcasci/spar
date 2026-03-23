package solution

import (
	"sort"
	"strings"
)

func SolveNQueens(n int) [][]string {
	var result [][]string
	cols := make(map[int]bool)
	diag1 := make(map[int]bool)
	diag2 := make(map[int]bool)
	board := make([][]byte, n)
	for i := range board {
		board[i] = make([]byte, n)
		for j := range board[i] {
			board[i][j] = '.'
		}
	}

	var backtrack func(row int)
	backtrack = func(row int) {
		if row == n {
			sol := make([]string, n)
			for i := range board {
				sol[i] = string(board[i])
			}
			result = append(result, sol)
			return
		}
		for c := 0; c < n; c++ {
			if cols[c] || diag1[row-c] || diag2[row+c] {
				continue
			}
			cols[c] = true
			diag1[row-c] = true
			diag2[row+c] = true
			board[row][c] = 'Q'
			backtrack(row + 1)
			board[row][c] = '.'
			delete(cols, c)
			delete(diag1, row-c)
			delete(diag2, row+c)
		}
	}

	backtrack(0)
	sort.Slice(result, func(i, j int) bool {
		return strings.Join(result[i], ",") < strings.Join(result[j], ",")
	})
	return result
}
