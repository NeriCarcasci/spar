function topKFrequent(nums, k) {
    const freq = new Map();
    for (const n of nums) {
        freq.set(n, (freq.get(n) || 0) + 1);
    }

    const buckets = Array.from({ length: nums.length + 1 }, () => []);
    for (const [num, count] of freq) {
        buckets[count].push(num);
    }

    const result = [];
    for (let i = buckets.length - 1; i > 0 && result.length < k; i--) {
        result.push(...buckets[i]);
    }
    return result.slice(0, k).sort((a, b) => a - b);
}

module.exports = { topKFrequent };
