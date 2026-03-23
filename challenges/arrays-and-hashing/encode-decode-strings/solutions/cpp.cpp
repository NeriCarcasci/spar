#include <vector>
#include <string>

class Solution {
public:
    std::string encode(std::vector<std::string>& strs) {
        std::string result;
        for (auto& s : strs) {
            result += std::to_string(s.size()) + "#" + s;
        }
        return result;
    }

    std::vector<std::string> decode(std::string s) {
        std::vector<std::string> result;
        size_t i = 0;
        while (i < s.size()) {
            size_t j = s.find('#', i);
            int length = std::stoi(s.substr(i, j - i));
            result.push_back(s.substr(j + 1, length));
            i = j + 1 + length;
        }
        return result;
    }
};
