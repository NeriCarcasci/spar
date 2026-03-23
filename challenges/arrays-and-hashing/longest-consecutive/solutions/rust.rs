use std::collections::HashSet;

pub fn longest_consecutive(nums: Vec<i32>) -> i32 {
    let num_set: HashSet<i32> = nums.into_iter().collect();
    let mut longest = 0;
    for &n in &num_set {
        if num_set.contains(&(n - 1)) {
            continue;
        }
        let mut length = 1;
        while num_set.contains(&(n + length)) {
            length += 1;
        }
        longest = longest.max(length);
    }
    longest
}
