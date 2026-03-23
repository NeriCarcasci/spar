#include <vector>
#include <algorithm>

class Solution {
public:
    int trap(std::vector<int>& height) {
        int left = 0, right = static_cast<int>(height.size()) - 1;
        int leftMax = 0, rightMax = 0, water = 0;
        while (left < right) {
            if (height[left] < height[right]) {
                leftMax = std::max(leftMax, height[left]);
                water += leftMax - height[left];
                left++;
            } else {
                rightMax = std::max(rightMax, height[right]);
                water += rightMax - height[right];
                right--;
            }
        }
        return water;
    }
};
