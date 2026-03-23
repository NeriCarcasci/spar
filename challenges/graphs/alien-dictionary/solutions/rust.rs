use std::collections::{HashMap, HashSet, VecDeque};

pub fn alien_order(words: Vec<String>) -> String {
    let mut adj: HashMap<u8, HashSet<u8>> = HashMap::new();
    let mut in_degree: HashMap<u8, i32> = HashMap::new();

    for word in &words {
        for &ch in word.as_bytes() {
            adj.entry(ch).or_insert_with(HashSet::new);
            in_degree.entry(ch).or_insert(0);
        }
    }

    for i in 0..words.len() - 1 {
        let w1 = words[i].as_bytes();
        let w2 = words[i + 1].as_bytes();
        let min_len = w1.len().min(w2.len());
        if w1.len() > w2.len() && w1[..min_len] == w2[..min_len] {
            return String::new();
        }
        for j in 0..min_len {
            if w1[j] != w2[j] {
                if adj.get_mut(&w1[j]).unwrap().insert(w2[j]) {
                    *in_degree.get_mut(&w2[j]).unwrap() += 1;
                }
                break;
            }
        }
    }

    let mut queue = VecDeque::new();
    for (&ch, &deg) in &in_degree {
        if deg == 0 {
            queue.push_back(ch);
        }
    }

    let mut result = Vec::new();
    while let Some(ch) = queue.pop_front() {
        result.push(ch);
        if let Some(neighbors) = adj.get(&ch) {
            for &neighbor in neighbors {
                let deg = in_degree.get_mut(&neighbor).unwrap();
                *deg -= 1;
                if *deg == 0 {
                    queue.push_back(neighbor);
                }
            }
        }
    }

    if result.len() != in_degree.len() {
        return String::new();
    }
    String::from_utf8(result).unwrap()
}
