#include <vector>
#include <algorithm>

class Solution {
public:
    std::vector<std::vector<int>> kClosest(std::vector<std::vector<int>>& points, int k) {
        std::sort(points.begin(), points.end(), [](const auto& a, const auto& b) {
            int da = a[0] * a[0] + a[1] * a[1];
            int db = b[0] * b[0] + b[1] * b[1];
            if (da != db) return da < db;
            if (a[0] != b[0]) return a[0] < b[0];
            return a[1] < b[1];
        });
        points.resize(k);
        return points;
    }
};
