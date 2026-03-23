pub fn valid_tree(n: i32, edges: Vec<Vec<i32>>) -> bool {
    let n = n as usize;
    if edges.len() != n - 1 {
        return false;
    }
    let mut parent: Vec<usize> = (0..n).collect();
    let mut rank = vec![0usize; n];

    fn find(parent: &mut Vec<usize>, mut x: usize) -> usize {
        while parent[x] != x {
            parent[x] = parent[parent[x]];
            x = parent[x];
        }
        x
    }

    for e in &edges {
        let mut px = find(&mut parent, e[0] as usize);
        let mut py = find(&mut parent, e[1] as usize);
        if px == py {
            return false;
        }
        if rank[px] < rank[py] {
            std::mem::swap(&mut px, &mut py);
        }
        parent[py] = px;
        if rank[px] == rank[py] {
            rank[px] += 1;
        }
    }
    true
}
