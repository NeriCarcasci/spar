def find_words(board: list[list[str]], words: list[str]) -> list[str]:
    root = {}
    for word in words:
        node = root
        for ch in word:
            node = node.setdefault(ch, {})
        node["#"] = word

    rows, cols = len(board), len(board[0])
    result = []

    def dfs(r, c, node):
        ch = board[r][c]
        if ch not in node:
            return
        next_node = node[ch]
        if "#" in next_node:
            result.append(next_node.pop("#"))
        board[r][c] = "."
        for dr, dc in ((0, 1), (0, -1), (1, 0), (-1, 0)):
            nr, nc = r + dr, c + dc
            if 0 <= nr < rows and 0 <= nc < cols and board[nr][nc] != ".":
                dfs(nr, nc, next_node)
        board[r][c] = ch
        if not next_node:
            del node[ch]

    for r in range(rows):
        for c in range(cols):
            dfs(r, c, root)

    result.sort()
    return result
