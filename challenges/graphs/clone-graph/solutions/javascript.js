class Node {
    constructor(val, neighbors) {
        this.val = val === undefined ? 0 : val;
        this.neighbors = neighbors === undefined ? [] : neighbors;
    }
}

function cloneGraph(node) {
    if (!node) return null;
    const cloned = new Map();

    function dfs(n) {
        if (cloned.has(n)) return cloned.get(n);
        const copy = new Node(n.val);
        cloned.set(n, copy);
        for (const neighbor of n.neighbors) {
            copy.neighbors.push(dfs(neighbor));
        }
        return copy;
    }

    return dfs(node);
}

module.exports = { Node, cloneGraph };
