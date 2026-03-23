function findWords(board, words) {
    const root = {};
    for (const word of words) {
        let node = root;
        for (const ch of word) {
            if (!node[ch]) node[ch] = {};
            node = node[ch];
        }
        node.word = word;
    }

    const rows = board.length;
    const cols = board[0].length;
    const result = [];
    const dirs = [[0, 1], [0, -1], [1, 0], [-1, 0]];

    function dfs(r, c, node) {
        const ch = board[r][c];
        if (!node[ch]) return;
        const next = node[ch];
        if (next.word) {
            result.push(next.word);
            delete next.word;
        }
        board[r][c] = "#";
        for (const [dr, dc] of dirs) {
            const nr = r + dr, nc = c + dc;
            if (nr >= 0 && nr < rows && nc >= 0 && nc < cols && board[nr][nc] !== "#") {
                dfs(nr, nc, next);
            }
        }
        board[r][c] = ch;
    }

    for (let r = 0; r < rows; r++) {
        for (let c = 0; c < cols; c++) {
            dfs(r, c, root);
        }
    }

    return result.sort();
}

module.exports = { findWords };
