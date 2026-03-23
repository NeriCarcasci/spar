from collections import defaultdict

def group_anagrams(strs: list[str]) -> list[list[str]]:
    groups = defaultdict(list)
    for s in strs:
        key = tuple(sorted(s))
        groups[key].append(s)
    result = [sorted(g) for g in groups.values()]
    result.sort(key=lambda g: g[0])
    return result
