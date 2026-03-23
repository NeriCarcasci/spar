use std::collections::HashMap;

pub fn top_k_frequent(nums: Vec<i32>, k: i32) -> Vec<i32> {
    let mut freq = HashMap::new();
    for &n in &nums {
        *freq.entry(n).or_insert(0) += 1;
    }

    let mut buckets = vec![vec![]; nums.len() + 1];
    for (num, count) in freq {
        buckets[count as usize].push(num);
    }

    let mut result = Vec::with_capacity(k as usize);
    for bucket in buckets.iter().rev() {
        for &num in bucket {
            result.push(num);
            if result.len() == k as usize {
                result.sort();
                return result;
            }
        }
    }
    result.sort();
    result
}
