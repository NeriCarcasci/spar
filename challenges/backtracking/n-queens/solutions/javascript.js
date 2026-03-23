function solveNQueens(n) {
    const result = [];
    const cols = new Set();
    const diag1 = new Set();
    const diag2 = new Set();
    const board = Array.from({length: n}, () => Array(n).fill("."));

    function backtrack(row) {
        if (row === n) {
            result.push(board.map(r => r.join("")));
            return;
        }
        for (let c = 0; c < n; c++) {
            if (cols.has(c) || diag1.has(row - c) || diag2.has(row + c)) continue;
            cols.add(c);
            diag1.add(row - c);
            diag2.add(row + c);
            board[row][c] = "Q";
            backtrack(row + 1);
            board[row][c] = ".";
            cols.delete(c);
            diag1.delete(row - c);
            diag2.delete(row + c);
        }
    }

    backtrack(0);
    result.sort((a, b) => a.join(",").localeCompare(b.join(",")));
    return result;
}

module.exports = { solveNQueens };
