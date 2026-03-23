#include <vector>
#include <algorithm>
#include <climits>

class Solution {
public:
    int maxProfit(std::vector<int>& prices) {
        int minPrice = INT_MAX, best = 0;
        for (int p : prices) {
            minPrice = std::min(minPrice, p);
            best = std::max(best, p - minPrice);
        }
        return best;
    }
};
