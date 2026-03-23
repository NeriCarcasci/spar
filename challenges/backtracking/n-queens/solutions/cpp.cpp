#include <vector>
#include <string>
#include <algorithm>
#include <unordered_set>

class Solution {
    int n;
    std::unordered_set<int> cols, diag1, diag2;
    std::vector<std::string> board;
    std::vector<std::vector<std::string>> result;

    void backtrack(int row) {
        if (row == n) {
            result.push_back(board);
            return;
        }
        for (int c = 0; c < n; c++) {
            if (cols.count(c) || diag1.count(row - c) || diag2.count(row + c)) continue;
            cols.insert(c);
            diag1.insert(row - c);
            diag2.insert(row + c);
            board[row][c] = 'Q';
            backtrack(row + 1);
            board[row][c] = '.';
            cols.erase(c);
            diag1.erase(row - c);
            diag2.erase(row + c);
        }
    }

public:
    std::vector<std::vector<std::string>> solveNQueens(int n) {
        this->n = n;
        board.assign(n, std::string(n, '.'));
        backtrack(0);
        std::sort(result.begin(), result.end());
        return result;
    }
};
