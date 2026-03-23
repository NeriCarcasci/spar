#include <string>
#include <unordered_map>
#include <algorithm>

class Solution {
public:
    int lengthOfLongestSubstring(std::string s) {
        std::unordered_map<char, int> lastSeen;
        int start = 0, longest = 0;
        for (int i = 0; i < static_cast<int>(s.size()); i++) {
            auto it = lastSeen.find(s[i]);
            if (it != lastSeen.end() && it->second >= start) {
                start = it->second + 1;
            }
            lastSeen[s[i]] = i;
            longest = std::max(longest, i - start + 1);
        }
        return longest;
    }
};
