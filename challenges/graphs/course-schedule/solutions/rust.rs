use std::collections::VecDeque;

pub fn can_finish(num_courses: i32, prerequisites: Vec<Vec<i32>>) -> bool {
    let n = num_courses as usize;
    let mut graph = vec![vec![]; n];
    let mut in_degree = vec![0; n];

    for p in &prerequisites {
        graph[p[1] as usize].push(p[0] as usize);
        in_degree[p[0] as usize] += 1;
    }

    let mut queue = VecDeque::new();
    for i in 0..n {
        if in_degree[i] == 0 {
            queue.push_back(i);
        }
    }

    let mut count = 0;
    while let Some(node) = queue.pop_front() {
        count += 1;
        for &neighbor in &graph[node] {
            in_degree[neighbor] -= 1;
            if in_degree[neighbor] == 0 {
                queue.push_back(neighbor);
            }
        }
    }
    count == n
}
