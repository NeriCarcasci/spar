pub fn least_interval(tasks: Vec<char>, n: i32) -> i32 {
    let mut freq = [0i32; 26];
    for &t in &tasks {
        freq[(t as u8 - b'A') as usize] += 1;
    }
    let max_freq = *freq.iter().max().unwrap();
    let count_max = freq.iter().filter(|&&f| f == max_freq).count() as i32;
    std::cmp::max(tasks.len() as i32, (max_freq - 1) * (n + 1) + count_max)
}
