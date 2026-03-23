function maxSlidingWindow(nums, k) {
    const dq = [];
    const result = [];
    for (let i = 0; i < nums.length; i++) {
        while (dq.length && dq[0] < i - k + 1) dq.shift();
        while (dq.length && nums[dq[dq.length - 1]] <= nums[i]) dq.pop();
        dq.push(i);
        if (i >= k - 1) result.push(nums[dq[0]]);
    }
    return result;
}
module.exports = { maxSlidingWindow };
