class TreeNode {
    constructor(val = 0, left = null, right = null) {
        this.val = val;
        this.left = left;
        this.right = right;
    }
}

function kthSmallest(root, k) {
    const stack = [];
    let cur = root;
    let count = 0;
    while (cur || stack.length > 0) {
        while (cur) {
            stack.push(cur);
            cur = cur.left;
        }
        cur = stack.pop();
        count++;
        if (count === k) return cur.val;
        cur = cur.right;
    }
    return -1;
}

module.exports = { TreeNode, kthSmallest };
