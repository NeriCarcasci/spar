#include <vector>
#include <algorithm>
#include <climits>

class Solution {
public:
    double findMedianSortedArrays(std::vector<int>& nums1, std::vector<int>& nums2) {
        if (nums1.size() > nums2.size()) std::swap(nums1, nums2);
        int m = static_cast<int>(nums1.size()), n = static_cast<int>(nums2.size());
        int left = 0, right = m;
        while (left <= right) {
            int i = (left + right) / 2;
            int j = (m + n + 1) / 2 - i;
            int left1 = i > 0 ? nums1[i - 1] : INT_MIN;
            int right1 = i < m ? nums1[i] : INT_MAX;
            int left2 = j > 0 ? nums2[j - 1] : INT_MIN;
            int right2 = j < n ? nums2[j] : INT_MAX;
            if (left1 <= right2 && left2 <= right1) {
                if ((m + n) % 2 == 0) return (std::max(left1, left2) + std::min(right1, right2)) / 2.0;
                return std::max(left1, left2);
            } else if (left1 > right2) { right = i - 1; }
            else { left = i + 1; }
        }
        return 0;
    }
};
