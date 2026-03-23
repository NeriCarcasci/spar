class Node:
    def __init__(self, val=0, neighbors=None):
        self.val = val
        self.neighbors = neighbors if neighbors is not None else []


def clone_graph(node: "Node | None") -> "Node | None":
    if not node:
        return None
    cloned = {}

    def dfs(n):
        if n in cloned:
            return cloned[n]
        copy = Node(n.val)
        cloned[n] = copy
        for neighbor in n.neighbors:
            copy.neighbors.append(dfs(neighbor))
        return copy

    return dfs(node)
