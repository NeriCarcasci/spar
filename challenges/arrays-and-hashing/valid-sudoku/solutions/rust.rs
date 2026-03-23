pub fn is_valid_sudoku(board: Vec<Vec<String>>) -> bool {
    let mut rows = [[false; 9]; 9];
    let mut cols = [[false; 9]; 9];
    let mut boxes = [[false; 9]; 9];
    for r in 0..9 {
        for c in 0..9 {
            if board[r][c] == "." {
                continue;
            }
            let d = (board[r][c].as_bytes()[0] - b'1') as usize;
            let b = (r / 3) * 3 + c / 3;
            if rows[r][d] || cols[c][d] || boxes[b][d] {
                return false;
            }
            rows[r][d] = true;
            cols[c][d] = true;
            boxes[b][d] = true;
        }
    }
    true
}
