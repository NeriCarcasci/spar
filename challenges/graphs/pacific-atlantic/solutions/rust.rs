pub fn pacific_atlantic(heights: Vec<Vec<i32>>) -> Vec<Vec<i32>> {
    if heights.is_empty() {
        return vec![];
    }
    let rows = heights.len();
    let cols = heights[0].len();
    let mut pacific = vec![vec![false; cols]; rows];
    let mut atlantic = vec![vec![false; cols]; rows];

    for c in 0..cols {
        dfs(&heights, &mut pacific, 0, c, rows, cols, heights[0][c]);
        dfs(&heights, &mut atlantic, rows - 1, c, rows, cols, heights[rows - 1][c]);
    }
    for r in 0..rows {
        dfs(&heights, &mut pacific, r, 0, rows, cols, heights[r][0]);
        dfs(&heights, &mut atlantic, r, cols - 1, rows, cols, heights[r][cols - 1]);
    }

    let mut result = vec![];
    for r in 0..rows {
        for c in 0..cols {
            if pacific[r][c] && atlantic[r][c] {
                result.push(vec![r as i32, c as i32]);
            }
        }
    }
    result
}

fn dfs(heights: &[Vec<i32>], visited: &mut Vec<Vec<bool>>, r: usize, c: usize, rows: usize, cols: usize, prev_height: i32) {
    if r >= rows || c >= cols || visited[r][c] || heights[r][c] < prev_height {
        return;
    }
    visited[r][c] = true;
    if r + 1 < rows { dfs(heights, visited, r + 1, c, rows, cols, heights[r][c]); }
    if r > 0 { dfs(heights, visited, r - 1, c, rows, cols, heights[r][c]); }
    if c + 1 < cols { dfs(heights, visited, r, c + 1, rows, cols, heights[r][c]); }
    if c > 0 { dfs(heights, visited, r, c - 1, rows, cols, heights[r][c]); }
}
