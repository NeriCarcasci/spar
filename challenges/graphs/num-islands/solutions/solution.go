package solution

func NumIslands(grid [][]string) int {
	if len(grid) == 0 {
		return 0
	}
	rows, cols := len(grid), len(grid[0])
	count := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] == "1" {
				count++
				dfs(grid, r, c, rows, cols)
			}
		}
	}
	return count
}

func dfs(grid [][]string, r, c, rows, cols int) {
	if r < 0 || r >= rows || c < 0 || c >= cols || grid[r][c] != "1" {
		return
	}
	grid[r][c] = "0"
	dfs(grid, r+1, c, rows, cols)
	dfs(grid, r-1, c, rows, cols)
	dfs(grid, r, c+1, rows, cols)
	dfs(grid, r, c-1, rows, cols)
}
