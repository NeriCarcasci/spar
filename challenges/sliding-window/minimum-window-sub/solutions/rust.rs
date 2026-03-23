use std::collections::HashMap;

pub fn min_window(s: String, t: String) -> String {
    let sb = s.as_bytes();
    let mut need: HashMap<u8, i32> = HashMap::new();
    for &b in t.as_bytes() {
        *need.entry(b).or_default() += 1;
    }
    let mut missing = need.len() as i32;
    let (mut left, mut best_left, mut best_len) = (0usize, 0usize, sb.len() + 1);
    for right in 0..sb.len() {
        let e = need.entry(sb[right]).or_default();
        *e -= 1;
        if *e == 0 { missing -= 1; }
        while missing == 0 {
            if right - left + 1 < best_len {
                best_left = left;
                best_len = right - left + 1;
            }
            let e = need.entry(sb[left]).or_default();
            *e += 1;
            if *e > 0 { missing += 1; }
            left += 1;
        }
    }
    if best_len > sb.len() { String::new() } else { s[best_left..best_left + best_len].to_string() }
}
