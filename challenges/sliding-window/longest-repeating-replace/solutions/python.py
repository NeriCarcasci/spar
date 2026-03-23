def character_replacement(s: str, k: int) -> int:
    counts = {}
    left = 0
    max_freq = 0
    longest = 0
    for right in range(len(s)):
        counts[s[right]] = counts.get(s[right], 0) + 1
        max_freq = max(max_freq, counts[s[right]])
        while (right - left + 1) - max_freq > k:
            counts[s[left]] -= 1
            left += 1
        longest = max(longest, right - left + 1)
    return longest
