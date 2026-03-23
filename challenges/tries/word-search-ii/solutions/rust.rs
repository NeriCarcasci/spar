use std::collections::HashMap;

struct TrieNode {
    children: HashMap<char, TrieNode>,
    word: Option<String>,
}

impl TrieNode {
    fn new() -> Self {
        TrieNode { children: HashMap::new(), word: None }
    }
}

pub fn find_words(mut board: Vec<Vec<char>>, words: Vec<String>) -> Vec<String> {
    let mut root = TrieNode::new();
    for word in &words {
        let mut node = &mut root;
        for ch in word.chars() {
            node = node.children.entry(ch).or_insert_with(TrieNode::new);
        }
        node.word = Some(word.clone());
    }

    let rows = board.len();
    let cols = board[0].len();
    let mut result = Vec::new();

    for r in 0..rows {
        for c in 0..cols {
            dfs(&mut board, r, c, &mut root, &mut result);
        }
    }
    result.sort();
    result
}

fn dfs(board: &mut Vec<Vec<char>>, r: usize, c: usize, node: &mut TrieNode, result: &mut Vec<String>) {
    let ch = board[r][c];
    if ch == '#' {
        return;
    }
    if let Some(next) = node.children.get_mut(&ch) {
        if let Some(word) = next.word.take() {
            result.push(word);
        }
        board[r][c] = '#';
        let dirs: [(i32, i32); 4] = [(0, 1), (0, -1), (1, 0), (-1, 0)];
        for (dr, dc) in dirs {
            let nr = r as i32 + dr;
            let nc = c as i32 + dc;
            if nr >= 0 && nr < board.len() as i32 && nc >= 0 && nc < board[0].len() as i32 {
                dfs(board, nr as usize, nc as usize, next, result);
            }
        }
        board[r][c] = ch;
    }
}
