pub fn is_palindrome(s: String) -> bool {
    let chars: Vec<u8> = s.bytes().filter(|b| b.is_ascii_alphanumeric()).map(|b| b.to_ascii_lowercase()).collect();
    let n = chars.len();
    for i in 0..n / 2 {
        if chars[i] != chars[n - 1 - i] {
            return false;
        }
    }
    true
}
