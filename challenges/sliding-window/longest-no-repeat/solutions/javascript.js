function lengthOfLongestSubstring(s) {
    const lastSeen = new Map();
    let start = 0, longest = 0;
    for (let i = 0; i < s.length; i++) {
        if (lastSeen.has(s[i]) && lastSeen.get(s[i]) >= start) {
            start = lastSeen.get(s[i]) + 1;
        }
        lastSeen.set(s[i], i);
        longest = Math.max(longest, i - start + 1);
    }
    return longest;
}
module.exports = { lengthOfLongestSubstring };
