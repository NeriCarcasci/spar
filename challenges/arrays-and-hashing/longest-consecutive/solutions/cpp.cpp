#include <vector>
#include <unordered_set>
#include <algorithm>

class Solution {
public:
    int longestConsecutive(std::vector<int>& nums) {
        std::unordered_set<int> numSet(nums.begin(), nums.end());
        int longest = 0;
        for (int n : numSet) {
            if (numSet.count(n - 1)) continue;
            int length = 1;
            while (numSet.count(n + length)) length++;
            longest = std::max(longest, length);
        }
        return longest;
    }
};
