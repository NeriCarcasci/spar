#include <vector>
#include <stack>
#include <algorithm>

class Solution {
public:
    int largestRectangleArea(std::vector<int>& heights) {
        std::stack<std::pair<int, int>> stk;
        int maxArea = 0;
        for (int i = 0; i < static_cast<int>(heights.size()); i++) {
            int start = i;
            while (!stk.empty() && stk.top().second > heights[i]) {
                auto [idx, h] = stk.top(); stk.pop();
                maxArea = std::max(maxArea, h * (i - idx));
                start = idx;
            }
            stk.push({start, heights[i]});
        }
        while (!stk.empty()) {
            auto [idx, h] = stk.top(); stk.pop();
            maxArea = std::max(maxArea, h * (static_cast<int>(heights.size()) - idx));
        }
        return maxArea;
    }
};
