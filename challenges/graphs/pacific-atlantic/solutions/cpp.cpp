#include <vector>

class Solution {
public:
    std::vector<std::vector<int>> pacificAtlantic(std::vector<std::vector<int>>& heights) {
        if (heights.empty()) return {};
        int rows = static_cast<int>(heights.size());
        int cols = static_cast<int>(heights[0].size());
        std::vector<std::vector<bool>> pacific(rows, std::vector<bool>(cols, false));
        std::vector<std::vector<bool>> atlantic(rows, std::vector<bool>(cols, false));

        for (int c = 0; c < cols; c++) {
            dfs(heights, pacific, 0, c, rows, cols, heights[0][c]);
            dfs(heights, atlantic, rows - 1, c, rows, cols, heights[rows - 1][c]);
        }
        for (int r = 0; r < rows; r++) {
            dfs(heights, pacific, r, 0, rows, cols, heights[r][0]);
            dfs(heights, atlantic, r, cols - 1, rows, cols, heights[r][cols - 1]);
        }

        std::vector<std::vector<int>> result;
        for (int r = 0; r < rows; r++) {
            for (int c = 0; c < cols; c++) {
                if (pacific[r][c] && atlantic[r][c]) {
                    result.push_back({r, c});
                }
            }
        }
        return result;
    }

    void dfs(std::vector<std::vector<int>>& heights, std::vector<std::vector<bool>>& visited,
             int r, int c, int rows, int cols, int prevHeight) {
        if (r < 0 || r >= rows || c < 0 || c >= cols) return;
        if (visited[r][c] || heights[r][c] < prevHeight) return;
        visited[r][c] = true;
        dfs(heights, visited, r + 1, c, rows, cols, heights[r][c]);
        dfs(heights, visited, r - 1, c, rows, cols, heights[r][c]);
        dfs(heights, visited, r, c + 1, rows, cols, heights[r][c]);
        dfs(heights, visited, r, c - 1, rows, cols, heights[r][c]);
    }
};
