package solution

func IsValidSudoku(board [][]string) bool {
	var rows, cols, boxes [9][9]bool
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if board[r][c] == "." {
				continue
			}
			d := board[r][c][0] - '1'
			box := (r/3)*3 + c/3
			if rows[r][d] || cols[c][d] || boxes[box][d] {
				return false
			}
			rows[r][d] = true
			cols[c][d] = true
			boxes[box][d] = true
		}
	}
	return true
}
