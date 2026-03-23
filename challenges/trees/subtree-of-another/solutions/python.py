class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right

def is_subtree(root: TreeNode, subRoot: TreeNode) -> bool:
    def is_same(p, q):
        if not p and not q:
            return True
        if not p or not q:
            return False
        return p.val == q.val and is_same(p.left, q.left) and is_same(p.right, q.right)

    if not root:
        return False
    if is_same(root, subRoot):
        return True
    return is_subtree(root.left, subRoot) or is_subtree(root.right, subRoot)
