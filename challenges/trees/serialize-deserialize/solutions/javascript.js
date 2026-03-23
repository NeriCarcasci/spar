class TreeNode {
    constructor(val = 0, left = null, right = null) {
        this.val = val;
        this.left = left;
        this.right = right;
    }
}

function serialize(root) {
    const result = [];
    function dfs(node) {
        if (!node) {
            result.push("N");
            return;
        }
        result.push(String(node.val));
        dfs(node.left);
        dfs(node.right);
    }
    dfs(root);
    return result.join(",");
}

function deserialize(data) {
    const vals = data.split(",");
    let idx = 0;
    function dfs() {
        if (vals[idx] === "N") {
            idx++;
            return null;
        }
        const node = new TreeNode(parseInt(vals[idx]));
        idx++;
        node.left = dfs();
        node.right = dfs();
        return node;
    }
    return dfs();
}

module.exports = { TreeNode, serialize, deserialize };
