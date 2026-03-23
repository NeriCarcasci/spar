pub fn num_islands(mut grid: Vec<Vec<char>>) -> i32 {
    if grid.is_empty() {
        return 0;
    }
    let (rows, cols) = (grid.len(), grid[0].len());
    let mut count = 0;
    for r in 0..rows {
        for c in 0..cols {
            if grid[r][c] == '1' {
                count += 1;
                dfs(&mut grid, r, c, rows, cols);
            }
        }
    }
    count
}

fn dfs(grid: &mut Vec<Vec<char>>, r: usize, c: usize, rows: usize, cols: usize) {
    if r >= rows || c >= cols || grid[r][c] != '1' {
        return;
    }
    grid[r][c] = '0';
    if r + 1 < rows { dfs(grid, r + 1, c, rows, cols); }
    if r > 0 { dfs(grid, r - 1, c, rows, cols); }
    if c + 1 < cols { dfs(grid, r, c + 1, rows, cols); }
    if c > 0 { dfs(grid, r, c - 1, rows, cols); }
}
