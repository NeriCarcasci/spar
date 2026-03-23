#include <vector>
#include <unordered_map>
#include <algorithm>

class Solution {
public:
    std::vector<int> topKFrequent(std::vector<int>& nums, int k) {
        std::unordered_map<int, int> freq;
        for (int n : nums) {
            freq[n]++;
        }

        std::vector<std::vector<int>> buckets(nums.size() + 1);
        for (auto& [num, count] : freq) {
            buckets[count].push_back(num);
        }

        std::vector<int> result;
        for (int i = static_cast<int>(buckets.size()) - 1; i > 0 && static_cast<int>(result.size()) < k; i--) {
            for (int num : buckets[i]) {
                result.push_back(num);
            }
        }
        result.resize(k);
        std::sort(result.begin(), result.end());
        return result;
    }
};
