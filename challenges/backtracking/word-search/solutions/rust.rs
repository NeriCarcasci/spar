pub fn exist(mut board: Vec<Vec<char>>, word: String) -> bool {
    let word: Vec<char> = word.chars().collect();
    let rows = board.len();
    let cols = board[0].len();
    for r in 0..rows {
        for c in 0..cols {
            if dfs(&mut board, &word, r, c, 0) {
                return true;
            }
        }
    }
    false
}

fn dfs(board: &mut Vec<Vec<char>>, word: &[char], r: usize, c: usize, idx: usize) -> bool {
    if idx == word.len() {
        return true;
    }
    if r >= board.len() || c >= board[0].len() || board[r][c] != word[idx] {
        return false;
    }
    let tmp = board[r][c];
    board[r][c] = '#';
    let dirs: [(i32, i32); 4] = [(0, 1), (0, -1), (1, 0), (-1, 0)];
    for (dr, dc) in dirs {
        let nr = r as i32 + dr;
        let nc = c as i32 + dc;
        if nr >= 0 && nc >= 0 && dfs(board, word, nr as usize, nc as usize, idx + 1) {
            return true;
        }
    }
    board[r][c] = tmp;
    false
}
