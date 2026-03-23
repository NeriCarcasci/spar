def exist(board: list[list[str]], word: str) -> bool:
    rows, cols = len(board), len(board[0])

    def dfs(r, c, idx):
        if idx == len(word):
            return True
        if r < 0 or r >= rows or c < 0 or c >= cols or board[r][c] != word[idx]:
            return False
        tmp = board[r][c]
        board[r][c] = "#"
        for dr, dc in ((0, 1), (0, -1), (1, 0), (-1, 0)):
            if dfs(r + dr, c + dc, idx + 1):
                return True
        board[r][c] = tmp
        return False

    for r in range(rows):
        for c in range(cols):
            if dfs(r, c, 0):
                return True
    return False
