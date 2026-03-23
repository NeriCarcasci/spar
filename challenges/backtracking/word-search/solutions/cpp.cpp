#include <vector>
#include <string>

class Solution {
    int rows, cols;
    int dirs[4][2] = {{0,1},{0,-1},{1,0},{-1,0}};

    bool dfs(std::vector<std::vector<char>>& board, const std::string& word, int r, int c, int idx) {
        if (idx == (int)word.size()) return true;
        if (r < 0 || r >= rows || c < 0 || c >= cols || board[r][c] != word[idx]) return false;
        char tmp = board[r][c];
        board[r][c] = '#';
        for (auto& d : dirs) {
            if (dfs(board, word, r + d[0], c + d[1], idx + 1)) return true;
        }
        board[r][c] = tmp;
        return false;
    }

public:
    bool exist(std::vector<std::vector<char>>& board, std::string word) {
        rows = board.size();
        cols = board[0].size();
        for (int r = 0; r < rows; r++) {
            for (int c = 0; c < cols; c++) {
                if (dfs(board, word, r, c, 0)) return true;
            }
        }
        return false;
    }
};
