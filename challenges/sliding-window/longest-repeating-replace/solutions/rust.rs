pub fn character_replacement(s: String, k: i32) -> i32 {
    let bytes = s.as_bytes();
    let mut counts = [0i32; 26];
    let (mut left, mut max_freq, mut longest) = (0usize, 0i32, 0i32);
    for right in 0..bytes.len() {
        counts[(bytes[right] - b'A') as usize] += 1;
        max_freq = max_freq.max(counts[(bytes[right] - b'A') as usize]);
        while (right as i32 - left as i32 + 1) - max_freq > k {
            counts[(bytes[left] - b'A') as usize] -= 1;
            left += 1;
        }
        longest = longest.max(right as i32 - left as i32 + 1);
    }
    longest
}
