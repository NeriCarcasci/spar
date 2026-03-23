use std::collections::HashMap;

pub fn clone_graph(adj_list: &HashMap<i32, Vec<i32>>) -> HashMap<i32, Vec<i32>> {
    let mut cloned = HashMap::new();
    for (key, neighbors) in adj_list {
        cloned.insert(*key, neighbors.clone());
    }
    cloned
}
