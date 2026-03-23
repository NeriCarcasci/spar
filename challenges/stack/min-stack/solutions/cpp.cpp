#include <vector>
#include <algorithm>

class MinStack {
    std::vector<std::pair<int, int>> stack;
public:
    MinStack() {}

    void push(int val) {
        int currentMin = stack.empty() ? val : std::min(val, stack.back().second);
        stack.emplace_back(val, currentMin);
    }

    void pop() {
        stack.pop_back();
    }

    int top() {
        return stack.back().first;
    }

    int getMin() {
        return stack.back().second;
    }
};
