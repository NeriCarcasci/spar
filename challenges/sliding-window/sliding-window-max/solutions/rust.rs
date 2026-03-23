use std::collections::VecDeque;

pub fn max_sliding_window(nums: Vec<i32>, k: i32) -> Vec<i32> {
    let k = k as usize;
    let mut dq: VecDeque<usize> = VecDeque::new();
    let mut result = Vec::with_capacity(nums.len().saturating_sub(k) + 1);
    for i in 0..nums.len() {
        while dq.front().map_or(false, |&f| f + k <= i) {
            dq.pop_front();
        }
        while dq.back().map_or(false, |&b| nums[b] <= nums[i]) {
            dq.pop_back();
        }
        dq.push_back(i);
        if i >= k - 1 {
            result.push(nums[*dq.front().unwrap()]);
        }
    }
    result
}
