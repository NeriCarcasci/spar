pub fn daily_temperatures(temperatures: Vec<i32>) -> Vec<i32> {
    let n = temperatures.len();
    let mut answer = vec![0i32; n];
    let mut stack: Vec<usize> = Vec::new();
    for i in 0..n {
        while stack.last().map_or(false, |&j| temperatures[i] > temperatures[j]) {
            let j = stack.pop().unwrap();
            answer[j] = (i - j) as i32;
        }
        stack.push(i);
    }
    answer
}
