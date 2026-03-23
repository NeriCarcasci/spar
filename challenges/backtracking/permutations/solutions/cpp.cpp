#include <vector>
#include <algorithm>

class Solution {
    void backtrack(std::vector<int>& nums, std::vector<bool>& used,
                   std::vector<int>& path, std::vector<std::vector<int>>& result) {
        if (path.size() == nums.size()) {
            result.push_back(path);
            return;
        }
        for (int i = 0; i < (int)nums.size(); i++) {
            if (used[i]) continue;
            used[i] = true;
            path.push_back(nums[i]);
            backtrack(nums, used, path, result);
            path.pop_back();
            used[i] = false;
        }
    }
public:
    std::vector<std::vector<int>> permute(std::vector<int>& nums) {
        std::sort(nums.begin(), nums.end());
        std::vector<std::vector<int>> result;
        std::vector<int> path;
        std::vector<bool> used(nums.size(), false);
        backtrack(nums, used, path, result);
        return result;
    }
};
