function leastInterval(tasks, n) {
    const freq = {};
    for (const t of tasks) freq[t] = (freq[t] || 0) + 1;
    const maxFreq = Math.max(...Object.values(freq));
    const countMax = Object.values(freq).filter(v => v === maxFreq).length;
    return Math.max(tasks.length, (maxFreq - 1) * (n + 1) + countMax);
}

module.exports = { leastInterval };
