#include <vector>
#include <algorithm>

class Solution {
    void backtrack(std::vector<int>& candidates, int target, int start,
                   std::vector<int>& path, std::vector<std::vector<int>>& result) {
        if (target == 0) {
            result.push_back(path);
            return;
        }
        for (int i = start; i < (int)candidates.size(); i++) {
            if (candidates[i] > target) break;
            path.push_back(candidates[i]);
            backtrack(candidates, target - candidates[i], i, path, result);
            path.pop_back();
        }
    }
public:
    std::vector<std::vector<int>> combinationSum(std::vector<int>& candidates, int target) {
        std::sort(candidates.begin(), candidates.end());
        std::vector<std::vector<int>> result;
        std::vector<int> path;
        backtrack(candidates, target, 0, path, result);
        return result;
    }
};
