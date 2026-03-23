#include <vector>
#include <unordered_set>

class Solution {
public:
    bool containsDuplicate(std::vector<int>& nums) {
        std::unordered_set<int> seen;
        for (int n : nums) {
            if (!seen.insert(n).second) {
                return true;
            }
        }
        return false;
    }
};
