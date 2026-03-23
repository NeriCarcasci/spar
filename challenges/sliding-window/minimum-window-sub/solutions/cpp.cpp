#include <string>
#include <unordered_map>

class Solution {
public:
    std::string minWindow(std::string s, std::string t) {
        std::unordered_map<char, int> need;
        for (char c : t) need[c]++;
        int missing = static_cast<int>(need.size());
        int left = 0, bestLeft = 0, bestLen = static_cast<int>(s.size()) + 1;
        for (int right = 0; right < static_cast<int>(s.size()); right++) {
            if (--need[s[right]] == 0) missing--;
            while (missing == 0) {
                if (right - left + 1 < bestLen) {
                    bestLeft = left;
                    bestLen = right - left + 1;
                }
                if (++need[s[left]] > 0) missing++;
                left++;
            }
        }
        return bestLen > static_cast<int>(s.size()) ? "" : s.substr(bestLeft, bestLen);
    }
};
