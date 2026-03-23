function countComponents(n, edges) {
    const parent = Array.from({ length: n }, (_, i) => i);
    const rank = new Array(n).fill(0);

    function find(x) {
        while (parent[x] !== x) {
            parent[x] = parent[parent[x]];
            x = parent[x];
        }
        return x;
    }

    let components = n;
    for (const [a, b] of edges) {
        let px = find(a), py = find(b);
        if (px !== py) {
            if (rank[px] < rank[py]) [px, py] = [py, px];
            parent[py] = px;
            if (rank[px] === rank[py]) rank[px]++;
            components--;
        }
    }
    return components;
}

module.exports = { countComponents };
