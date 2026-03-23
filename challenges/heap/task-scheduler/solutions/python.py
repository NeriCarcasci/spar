from collections import Counter

def least_interval(tasks: list[str], n: int) -> int:
    freq = Counter(tasks)
    max_freq = max(freq.values())
    count_max = sum(1 for v in freq.values() if v == max_freq)
    return max(len(tasks), (max_freq - 1) * (n + 1) + count_max)
