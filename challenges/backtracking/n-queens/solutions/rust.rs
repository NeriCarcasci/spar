use std::collections::HashSet;

pub fn solve_n_queens(n: i32) -> Vec<Vec<String>> {
    let n = n as usize;
    let mut result = Vec::new();
    let mut cols = HashSet::new();
    let mut diag1 = HashSet::new();
    let mut diag2 = HashSet::new();
    let mut board = vec![vec!['.'; n]; n];
    backtrack(0, n, &mut cols, &mut diag1, &mut diag2, &mut board, &mut result);
    result.sort();
    result
}

fn backtrack(
    row: usize,
    n: usize,
    cols: &mut HashSet<usize>,
    diag1: &mut HashSet<i32>,
    diag2: &mut HashSet<usize>,
    board: &mut Vec<Vec<char>>,
    result: &mut Vec<Vec<String>>,
) {
    if row == n {
        result.push(board.iter().map(|r| r.iter().collect()).collect());
        return;
    }
    for c in 0..n {
        let d1 = row as i32 - c as i32;
        let d2 = row + c;
        if cols.contains(&c) || diag1.contains(&d1) || diag2.contains(&d2) {
            continue;
        }
        cols.insert(c);
        diag1.insert(d1);
        diag2.insert(d2);
        board[row][c] = 'Q';
        backtrack(row + 1, n, cols, diag1, diag2, board, result);
        board[row][c] = '.';
        cols.remove(&c);
        diag1.remove(&d1);
        diag2.remove(&d2);
    }
}
