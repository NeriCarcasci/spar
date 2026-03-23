class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right

def max_path_sum(root: TreeNode) -> int:
    result = [float('-inf')]

    def dfs(node):
        if not node:
            return 0
        left_gain = max(dfs(node.left), 0)
        right_gain = max(dfs(node.right), 0)
        result[0] = max(result[0], node.val + left_gain + right_gain)
        return node.val + max(left_gain, right_gain)

    dfs(root)
    return result[0]
