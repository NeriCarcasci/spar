class TreeNode {
    constructor(val = 0, left = null, right = null) {
        this.val = val;
        this.left = left;
        this.right = right;
    }
}

function lowestCommonAncestor(root, p, q) {
    let cur = root;
    while (cur) {
        if (p < cur.val && q < cur.val) {
            cur = cur.left;
        } else if (p > cur.val && q > cur.val) {
            cur = cur.right;
        } else {
            return cur.val;
        }
    }
    return -1;
}

module.exports = { TreeNode, lowestCommonAncestor };
