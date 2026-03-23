def max_area(height: list[int]) -> int:
    left, right = 0, len(height) - 1
    best = 0
    while left < right:
        width = right - left
        h = min(height[left], height[right])
        best = max(best, width * h)
        if height[left] < height[right]:
            left += 1
        else:
            right -= 1
    return best
