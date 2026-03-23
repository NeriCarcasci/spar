#include <vector>
#include <string>
#include <unordered_map>
#include <algorithm>

class Solution {
public:
    std::vector<std::vector<std::string>> groupAnagrams(std::vector<std::string>& strs) {
        std::unordered_map<std::string, std::vector<std::string>> groups;
        for (auto& s : strs) {
            auto key = s;
            std::sort(key.begin(), key.end());
            groups[key].push_back(s);
        }
        std::vector<std::vector<std::string>> result;
        result.reserve(groups.size());
        for (auto& [_, group] : groups) {
            std::sort(group.begin(), group.end());
            result.push_back(std::move(group));
        }
        std::sort(result.begin(), result.end(), [](auto& a, auto& b) {
            return a[0] < b[0];
        });
        return result;
    }
};
