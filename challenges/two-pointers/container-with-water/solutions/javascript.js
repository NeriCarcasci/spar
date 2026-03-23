function maxArea(height) {
    let left = 0, right = height.length - 1, best = 0;
    while (left < right) {
        const w = right - left;
        const h = Math.min(height[left], height[right]);
        best = Math.max(best, w * h);
        if (height[left] < height[right]) left++;
        else right--;
    }
    return best;
}
module.exports = { maxArea };
