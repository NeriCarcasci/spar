pub fn permute(mut nums: Vec<i32>) -> Vec<Vec<i32>> {
    nums.sort();
    let mut result = Vec::new();
    let mut used = vec![false; nums.len()];
    let mut path = Vec::new();
    backtrack(&nums, &mut used, &mut path, &mut result);
    result
}

fn backtrack(nums: &[i32], used: &mut Vec<bool>, path: &mut Vec<i32>, result: &mut Vec<Vec<i32>>) {
    if path.len() == nums.len() {
        result.push(path.clone());
        return;
    }
    for i in 0..nums.len() {
        if used[i] {
            continue;
        }
        used[i] = true;
        path.push(nums[i]);
        backtrack(nums, used, path, result);
        path.pop();
        used[i] = false;
    }
}
