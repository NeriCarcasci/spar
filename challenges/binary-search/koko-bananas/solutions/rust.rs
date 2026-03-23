pub fn min_eating_speed(piles: Vec<i32>, h: i32) -> i32 {
    let (mut left, mut right) = (1i64, *piles.iter().max().unwrap() as i64);
    while left < right {
        let mid = left + (right - left) / 2;
        let hours: i64 = piles.iter().map(|&p| (p as i64 + mid - 1) / mid).sum();
        if hours <= h as i64 { right = mid; }
        else { left = mid + 1; }
    }
    left as i32
}
