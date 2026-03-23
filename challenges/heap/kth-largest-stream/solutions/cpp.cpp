#include <vector>
#include <queue>

class KthLargest {
    int k;
    std::priority_queue<int, std::vector<int>, std::greater<int>> minHeap;
public:
    KthLargest(int k, std::vector<int>& nums) : k(k) {
        for (int n : nums) {
            minHeap.push(n);
            if ((int)minHeap.size() > k) minHeap.pop();
        }
    }

    int add(int val) {
        minHeap.push(val);
        if ((int)minHeap.size() > k) minHeap.pop();
        return minHeap.top();
    }
};
