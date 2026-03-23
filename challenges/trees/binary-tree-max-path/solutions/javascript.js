class TreeNode {
    constructor(val = 0, left = null, right = null) {
        this.val = val;
        this.left = left;
        this.right = right;
    }
}

function maxPathSum(root) {
    let result = -Infinity;

    function dfs(node) {
        if (!node) return 0;
        const leftGain = Math.max(dfs(node.left), 0);
        const rightGain = Math.max(dfs(node.right), 0);
        result = Math.max(result, node.val + leftGain + rightGain);
        return node.val + Math.max(leftGain, rightGain);
    }

    dfs(root);
    return result;
}

module.exports = { TreeNode, maxPathSum };
