from collections import deque

def alien_order(words: list[str]) -> str:
    adj = {}
    in_degree = {}

    for word in words:
        for ch in word:
            if ch not in in_degree:
                in_degree[ch] = 0
                adj[ch] = set()

    for i in range(len(words) - 1):
        w1, w2 = words[i], words[i + 1]
        min_len = min(len(w1), len(w2))
        if len(w1) > len(w2) and w1[:min_len] == w2[:min_len]:
            return ""
        for j in range(min_len):
            if w1[j] != w2[j]:
                if w2[j] not in adj[w1[j]]:
                    adj[w1[j]].add(w2[j])
                    in_degree[w2[j]] += 1
                break

    queue = deque()
    for ch in in_degree:
        if in_degree[ch] == 0:
            queue.append(ch)

    result = []
    while queue:
        ch = queue.popleft()
        result.append(ch)
        for neighbor in adj[ch]:
            in_degree[neighbor] -= 1
            if in_degree[neighbor] == 0:
                queue.append(neighbor)

    if len(result) != len(in_degree):
        return ""
    return "".join(result)
