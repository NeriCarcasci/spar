function groupAnagrams(strs) {
    const groups = new Map();
    for (const s of strs) {
        const key = [...s].sort().join("");
        if (!groups.has(key)) groups.set(key, []);
        groups.get(key).push(s);
    }
    const result = [...groups.values()].map(g => g.sort());
    result.sort((a, b) => a[0].localeCompare(b[0]));
    return result;
}

module.exports = { groupAnagrams };
