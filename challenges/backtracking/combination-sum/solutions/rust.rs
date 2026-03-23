pub fn combination_sum(mut candidates: Vec<i32>, target: i32) -> Vec<Vec<i32>> {
    candidates.sort();
    let mut result = Vec::new();
    let mut path = Vec::new();
    backtrack(&candidates, target, 0, &mut path, &mut result);
    result
}

fn backtrack(candidates: &[i32], remaining: i32, start: usize, path: &mut Vec<i32>, result: &mut Vec<Vec<i32>>) {
    if remaining == 0 {
        result.push(path.clone());
        return;
    }
    for i in start..candidates.len() {
        if candidates[i] > remaining {
            break;
        }
        path.push(candidates[i]);
        backtrack(candidates, remaining - candidates[i], i, path, result);
        path.pop();
    }
}
