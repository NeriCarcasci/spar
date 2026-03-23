function largestRectangleArea(heights) {
    const stack = [];
    let maxArea = 0;
    for (let i = 0; i < heights.length; i++) {
        let start = i;
        while (stack.length && stack[stack.length - 1][1] > heights[i]) {
            const [idx, h] = stack.pop();
            maxArea = Math.max(maxArea, h * (i - idx));
            start = idx;
        }
        stack.push([start, heights[i]]);
    }
    for (const [idx, h] of stack) {
        maxArea = Math.max(maxArea, h * (heights.length - idx));
    }
    return maxArea;
}
module.exports = { largestRectangleArea };
