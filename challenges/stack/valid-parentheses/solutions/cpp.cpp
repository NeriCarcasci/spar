#include <string>
#include <stack>
#include <unordered_map>

class Solution {
public:
    bool isValid(std::string s) {
        std::stack<char> stk;
        std::unordered_map<char, char> pairs = {{')', '('}, {'}', '{'}, {']', '['}};
        for (char c : s) {
            auto it = pairs.find(c);
            if (it != pairs.end()) {
                if (stk.empty() || stk.top() != it->second) return false;
                stk.pop();
            } else {
                stk.push(c);
            }
        }
        return stk.empty();
    }
};
