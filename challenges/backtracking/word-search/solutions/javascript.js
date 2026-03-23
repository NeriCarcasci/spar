function exist(board, word) {
    const rows = board.length;
    const cols = board[0].length;

    function dfs(r, c, idx) {
        if (idx === word.length) return true;
        if (r < 0 || r >= rows || c < 0 || c >= cols || board[r][c] !== word[idx]) return false;
        const tmp = board[r][c];
        board[r][c] = "#";
        const dirs = [[0,1],[0,-1],[1,0],[-1,0]];
        for (const [dr, dc] of dirs) {
            if (dfs(r + dr, c + dc, idx + 1)) return true;
        }
        board[r][c] = tmp;
        return false;
    }

    for (let r = 0; r < rows; r++) {
        for (let c = 0; c < cols; c++) {
            if (dfs(r, c, 0)) return true;
        }
    }
    return false;
}

module.exports = { exist };
