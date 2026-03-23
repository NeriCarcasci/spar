class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right

def serialize(root: TreeNode) -> str:
    result = []

    def dfs(node):
        if not node:
            result.append("N")
            return
        result.append(str(node.val))
        dfs(node.left)
        dfs(node.right)

    dfs(root)
    return ",".join(result)

def deserialize(data: str) -> TreeNode:
    vals = data.split(",")
    idx = [0]

    def dfs():
        if vals[idx[0]] == "N":
            idx[0] += 1
            return None
        node = TreeNode(int(vals[idx[0]]))
        idx[0] += 1
        node.left = dfs()
        node.right = dfs()
        return node

    return dfs()
