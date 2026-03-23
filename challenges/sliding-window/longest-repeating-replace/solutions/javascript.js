function characterReplacement(s, k) {
    const counts = new Array(26).fill(0);
    let left = 0, maxFreq = 0, longest = 0;
    for (let right = 0; right < s.length; right++) {
        counts[s.charCodeAt(right) - 65]++;
        maxFreq = Math.max(maxFreq, counts[s.charCodeAt(right) - 65]);
        while ((right - left + 1) - maxFreq > k) {
            counts[s.charCodeAt(left) - 65]--;
            left++;
        }
        longest = Math.max(longest, right - left + 1);
    }
    return longest;
}
module.exports = { characterReplacement };
