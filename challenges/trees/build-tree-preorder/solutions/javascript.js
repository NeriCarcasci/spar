class TreeNode {
    constructor(val = 0, left = null, right = null) {
        this.val = val;
        this.left = left;
        this.right = right;
    }
}

function buildTree(preorder, inorder) {
    if (preorder.length === 0) return null;
    const rootVal = preorder[0];
    const mid = inorder.indexOf(rootVal);
    const root = new TreeNode(rootVal);
    root.left = buildTree(preorder.slice(1, mid + 1), inorder.slice(0, mid));
    root.right = buildTree(preorder.slice(mid + 1), inorder.slice(mid + 1));
    return root;
}

module.exports = { TreeNode, buildTree };
