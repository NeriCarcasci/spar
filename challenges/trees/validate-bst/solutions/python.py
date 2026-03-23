class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right

def is_valid_bst(root: TreeNode) -> bool:
    def validate(node, lo, hi):
        if not node:
            return True
        if node.val <= lo or node.val >= hi:
            return False
        return validate(node.left, lo, node.val) and validate(node.right, node.val, hi)

    return validate(root, float('-inf'), float('inf'))
