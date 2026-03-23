pub fn max_area(height: Vec<i32>) -> i32 {
    let (mut left, mut right) = (0usize, height.len() - 1);
    let mut best = 0;
    while left < right {
        let w = (right - left) as i32;
        let h = height[left].min(height[right]);
        best = best.max(w * h);
        if height[left] < height[right] {
            left += 1;
        } else {
            right -= 1;
        }
    }
    best
}
