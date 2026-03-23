package solution

import "math"

func FindMedianSortedArrays(nums1 []int, nums2 []int) float64 {
	if len(nums1) > len(nums2) {
		nums1, nums2 = nums2, nums1
	}
	m, n := len(nums1), len(nums2)
	left, right := 0, m
	for left <= right {
		i := (left + right) / 2
		j := (m+n+1)/2 - i
		left1 := math.MinInt64
		if i > 0 { left1 = nums1[i-1] }
		right1 := math.MaxInt64
		if i < m { right1 = nums1[i] }
		left2 := math.MinInt64
		if j > 0 { left2 = nums2[j-1] }
		right2 := math.MaxInt64
		if j < n { right2 = nums2[j] }
		if left1 <= right2 && left2 <= right1 {
			if (m+n)%2 == 0 {
				return (float64(max(left1, left2)) + float64(min(right1, right2))) / 2.0
			}
			return float64(max(left1, left2))
		} else if left1 > right2 {
			right = i - 1
		} else {
			left = i + 1
		}
	}
	return 0
}
