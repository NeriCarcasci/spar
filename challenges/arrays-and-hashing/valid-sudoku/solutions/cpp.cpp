#include <vector>
#include <string>
#include <array>

class Solution {
public:
    bool isValidSudoku(std::vector<std::vector<std::string>>& board) {
        std::array<std::array<bool, 9>, 9> rows{}, cols{}, boxes{};
        for (int r = 0; r < 9; r++) {
            for (int c = 0; c < 9; c++) {
                if (board[r][c] == ".") continue;
                int d = board[r][c][0] - '1';
                int box = (r / 3) * 3 + c / 3;
                if (rows[r][d] || cols[c][d] || boxes[box][d]) return false;
                rows[r][d] = cols[c][d] = boxes[box][d] = true;
            }
        }
        return true;
    }
};
