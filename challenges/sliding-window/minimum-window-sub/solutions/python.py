from collections import Counter

def min_window(s: str, t: str) -> str:
    if not t or not s:
        return ""
    need = Counter(t)
    missing = len(need)
    left = 0
    best_left, best_right = 0, len(s)
    found = False
    for right, ch in enumerate(s):
        need[ch] -= 1
        if need[ch] == 0:
            missing -= 1
        while missing == 0:
            if right - left < best_right - best_left:
                best_left, best_right = left, right
                found = True
            need[s[left]] += 1
            if need[s[left]] > 0:
                missing += 1
            left += 1
    return s[best_left:best_right + 1] if found else ""
