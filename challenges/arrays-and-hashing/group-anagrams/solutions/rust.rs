use std::collections::HashMap;

pub fn group_anagrams(strs: Vec<String>) -> Vec<Vec<String>> {
    let mut groups: HashMap<Vec<u8>, Vec<String>> = HashMap::new();
    for s in strs {
        let mut key: Vec<u8> = s.bytes().collect();
        key.sort_unstable();
        groups.entry(key).or_default().push(s);
    }
    let mut result: Vec<Vec<String>> = groups.into_values().collect();
    for group in &mut result {
        group.sort();
    }
    result.sort_by(|a, b| a[0].cmp(&b[0]));
    result
}
