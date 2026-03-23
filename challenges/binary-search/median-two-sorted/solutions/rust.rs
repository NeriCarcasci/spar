pub fn find_median_sorted_arrays(nums1: Vec<i32>, nums2: Vec<i32>) -> f64 {
    let (a, b) = if nums1.len() <= nums2.len() { (&nums1, &nums2) } else { (&nums2, &nums1) };
    let (m, n) = (a.len(), b.len());
    let (mut lo, mut hi) = (0i32, m as i32);
    while lo <= hi {
        let i = ((lo + hi) / 2) as usize;
        let j = (m + n + 1) / 2 - i;
        let left1 = if i > 0 { a[i - 1] } else { i32::MIN };
        let right1 = if i < m { a[i] } else { i32::MAX };
        let left2 = if j > 0 { b[j - 1] } else { i32::MIN };
        let right2 = if j < n { b[j] } else { i32::MAX };
        if left1 <= right2 && left2 <= right1 {
            if (m + n) % 2 == 0 {
                return (left1.max(left2) as f64 + right1.min(right2) as f64) / 2.0;
            }
            return left1.max(left2) as f64;
        } else if left1 > right2 { hi = i as i32 - 1; }
        else { lo = i as i32 + 1; }
    }
    0.0
}
