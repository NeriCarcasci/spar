def solve_n_queens(n: int) -> list[list[str]]:
    result = []
    cols = set()
    diag1 = set()
    diag2 = set()
    board = [["." for _ in range(n)] for _ in range(n)]

    def backtrack(row):
        if row == n:
            result.append(["".join(r) for r in board])
            return
        for c in range(n):
            if c in cols or (row - c) in diag1 or (row + c) in diag2:
                continue
            cols.add(c)
            diag1.add(row - c)
            diag2.add(row + c)
            board[row][c] = "Q"
            backtrack(row + 1)
            board[row][c] = "."
            cols.remove(c)
            diag1.remove(row - c)
            diag2.remove(row + c)

    backtrack(0)
    result.sort()
    return result
