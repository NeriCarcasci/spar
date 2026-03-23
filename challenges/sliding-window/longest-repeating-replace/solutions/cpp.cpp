#include <string>
#include <array>
#include <algorithm>

class Solution {
public:
    int characterReplacement(std::string s, int k) {
        std::array<int, 26> counts{};
        int left = 0, maxFreq = 0, longest = 0;
        for (int right = 0; right < static_cast<int>(s.size()); right++) {
            counts[s[right] - 'A']++;
            maxFreq = std::max(maxFreq, counts[s[right] - 'A']);
            while ((right - left + 1) - maxFreq > k) {
                counts[s[left] - 'A']--;
                left++;
            }
            longest = std::max(longest, right - left + 1);
        }
        return longest;
    }
};
