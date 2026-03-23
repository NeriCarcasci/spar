#include <vector>
#include <algorithm>

class Solution {
    void backtrack(std::vector<int>& nums, int start,
                   std::vector<int>& path, std::vector<std::vector<int>>& result) {
        result.push_back(path);
        for (int i = start; i < (int)nums.size(); i++) {
            path.push_back(nums[i]);
            backtrack(nums, i + 1, path, result);
            path.pop_back();
        }
    }
public:
    std::vector<std::vector<int>> subsets(std::vector<int>& nums) {
        std::sort(nums.begin(), nums.end());
        std::vector<std::vector<int>> result;
        std::vector<int> path;
        backtrack(nums, 0, path, result);
        return result;
    }
};
