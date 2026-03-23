pub fn max_profit(prices: Vec<i32>) -> i32 {
    let mut min_price = i32::MAX;
    let mut best = 0;
    for &p in &prices {
        min_price = min_price.min(p);
        best = best.max(p - min_price);
    }
    best
}
