#include <queue>

class MedianFinder {
    std::priority_queue<int> lo;
    std::priority_queue<int, std::vector<int>, std::greater<int>> hi;
public:
    MedianFinder() {}

    void addNum(int num) {
        lo.push(num);
        hi.push(lo.top());
        lo.pop();
        if (hi.size() > lo.size()) {
            lo.push(hi.top());
            hi.pop();
        }
    }

    double findMedian() {
        if (lo.size() > hi.size()) {
            return lo.top();
        }
        return (lo.top() + hi.top()) / 2.0;
    }
};
