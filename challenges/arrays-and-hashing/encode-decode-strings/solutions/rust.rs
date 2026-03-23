pub fn encode(strs: Vec<String>) -> String {
    let mut result = String::new();
    for s in &strs {
        result.push_str(&s.len().to_string());
        result.push('#');
        result.push_str(s);
    }
    result
}

pub fn decode(s: String) -> Vec<String> {
    let mut result = Vec::new();
    let bytes = s.as_bytes();
    let mut i = 0;
    while i < bytes.len() {
        let j = s[i..].find('#').unwrap() + i;
        let length: usize = s[i..j].parse().unwrap();
        result.push(s[j + 1..j + 1 + length].to_string());
        i = j + 1 + length;
    }
    result
}
