package solution

func PacificAtlantic(heights [][]int) [][]int {
	if len(heights) == 0 {
		return nil
	}
	rows, cols := len(heights), len(heights[0])
	pacific := make([][]bool, rows)
	atlantic := make([][]bool, rows)
	for i := 0; i < rows; i++ {
		pacific[i] = make([]bool, cols)
		atlantic[i] = make([]bool, cols)
	}

	for c := 0; c < cols; c++ {
		dfsPa(heights, pacific, 0, c, rows, cols, heights[0][c])
		dfsPa(heights, atlantic, rows-1, c, rows, cols, heights[rows-1][c])
	}
	for r := 0; r < rows; r++ {
		dfsPa(heights, pacific, r, 0, rows, cols, heights[r][0])
		dfsPa(heights, atlantic, r, cols-1, rows, cols, heights[r][cols-1])
	}

	var result [][]int
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if pacific[r][c] && atlantic[r][c] {
				result = append(result, []int{r, c})
			}
		}
	}
	return result
}

func dfsPa(heights [][]int, visited [][]bool, r, c, rows, cols, prevHeight int) {
	if r < 0 || r >= rows || c < 0 || c >= cols || visited[r][c] || heights[r][c] < prevHeight {
		return
	}
	visited[r][c] = true
	dfsPa(heights, visited, r+1, c, rows, cols, heights[r][c])
	dfsPa(heights, visited, r-1, c, rows, cols, heights[r][c])
	dfsPa(heights, visited, r, c+1, rows, cols, heights[r][c])
	dfsPa(heights, visited, r, c-1, rows, cols, heights[r][c])
}
