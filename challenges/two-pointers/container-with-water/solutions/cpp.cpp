#include <vector>
#include <algorithm>

class Solution {
public:
    int maxArea(std::vector<int>& height) {
        int left = 0, right = static_cast<int>(height.size()) - 1, best = 0;
        while (left < right) {
            int w = right - left;
            int h = std::min(height[left], height[right]);
            best = std::max(best, w * h);
            if (height[left] < height[right]) left++;
            else right--;
        }
        return best;
    }
};
