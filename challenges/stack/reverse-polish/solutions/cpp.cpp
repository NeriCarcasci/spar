#include <vector>
#include <string>
#include <stack>

class Solution {
public:
    int evalRPN(std::vector<std::string>& tokens) {
        std::stack<int> stk;
        for (auto& t : tokens) {
            if (t == "+" || t == "-" || t == "*" || t == "/") {
                int b = stk.top(); stk.pop();
                int a = stk.top(); stk.pop();
                if (t == "+") stk.push(a + b);
                else if (t == "-") stk.push(a - b);
                else if (t == "*") stk.push(a * b);
                else stk.push(a / b);
            } else {
                stk.push(std::stoi(t));
            }
        }
        return stk.top();
    }
};
