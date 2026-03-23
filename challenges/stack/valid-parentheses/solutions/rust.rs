pub fn is_valid(s: String) -> bool {
    let mut stack = Vec::new();
    for b in s.bytes() {
        match b {
            b'(' | b'{' | b'[' => stack.push(b),
            b')' => if stack.pop() != Some(b'(') { return false; },
            b'}' => if stack.pop() != Some(b'{') { return false; },
            b']' => if stack.pop() != Some(b'[') { return false; },
            _ => {}
        }
    }
    stack.is_empty()
}
