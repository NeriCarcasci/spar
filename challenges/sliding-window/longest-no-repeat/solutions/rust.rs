use std::collections::HashMap;

pub fn length_of_longest_substring(s: String) -> i32 {
    let mut last_seen = HashMap::new();
    let mut start = 0i32;
    let mut longest = 0i32;
    for (i, b) in s.bytes().enumerate() {
        let i = i as i32;
        if let Some(&j) = last_seen.get(&b) {
            if j >= start {
                start = j + 1;
            }
        }
        last_seen.insert(b, i);
        longest = longest.max(i - start + 1);
    }
    longest
}
