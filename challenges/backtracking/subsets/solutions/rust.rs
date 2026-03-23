pub fn subsets(mut nums: Vec<i32>) -> Vec<Vec<i32>> {
    nums.sort();
    let mut result = Vec::new();
    let mut path = Vec::new();
    backtrack(&nums, 0, &mut path, &mut result);
    result
}

fn backtrack(nums: &[i32], start: usize, path: &mut Vec<i32>, result: &mut Vec<Vec<i32>>) {
    result.push(path.clone());
    for i in start..nums.len() {
        path.push(nums[i]);
        backtrack(nums, i + 1, path, result);
        path.pop();
    }
}
