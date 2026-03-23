function validTree(n, edges) {
    if (edges.length !== n - 1) return false;
    const parent = Array.from({ length: n }, (_, i) => i);
    const rank = new Array(n).fill(0);

    function find(x) {
        while (parent[x] !== x) {
            parent[x] = parent[parent[x]];
            x = parent[x];
        }
        return x;
    }

    for (const [a, b] of edges) {
        let px = find(a), py = find(b);
        if (px === py) return false;
        if (rank[px] < rank[py]) [px, py] = [py, px];
        parent[py] = px;
        if (rank[px] === rank[py]) rank[px]++;
    }
    return true;
}

module.exports = { validTree };
