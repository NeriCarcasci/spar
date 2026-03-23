#include <vector>
#include <algorithm>

class Solution {
public:
    int leastInterval(std::vector<char>& tasks, int n) {
        int freq[26] = {};
        for (char t : tasks) freq[t - 'A']++;
        int maxFreq = *std::max_element(freq, freq + 26);
        int countMax = 0;
        for (int f : freq) {
            if (f == maxFreq) countMax++;
        }
        return std::max((int)tasks.size(), (maxFreq - 1) * (n + 1) + countMax);
    }
};
