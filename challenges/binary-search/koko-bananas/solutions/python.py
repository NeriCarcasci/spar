import math

def min_eating_speed(piles: list[int], h: int) -> int:
    left, right = 1, max(piles)
    while left < right:
        mid = left + (right - left) // 2
        hours = sum(math.ceil(p / mid) for p in piles)
        if hours <= h:
            right = mid
        else:
            left = mid + 1
    return left
