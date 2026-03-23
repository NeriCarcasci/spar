function alienOrder(words) {
    const adj = new Map();
    const inDegree = new Map();

    for (const word of words) {
        for (const ch of word) {
            if (!inDegree.has(ch)) {
                inDegree.set(ch, 0);
                adj.set(ch, new Set());
            }
        }
    }

    for (let i = 0; i < words.length - 1; i++) {
        const w1 = words[i], w2 = words[i + 1];
        const minLen = Math.min(w1.length, w2.length);
        if (w1.length > w2.length && w1.slice(0, minLen) === w2.slice(0, minLen)) {
            return "";
        }
        for (let j = 0; j < minLen; j++) {
            if (w1[j] !== w2[j]) {
                if (!adj.get(w1[j]).has(w2[j])) {
                    adj.get(w1[j]).add(w2[j]);
                    inDegree.set(w2[j], inDegree.get(w2[j]) + 1);
                }
                break;
            }
        }
    }

    const queue = [];
    for (const [ch, deg] of inDegree) {
        if (deg === 0) queue.push(ch);
    }

    const result = [];
    while (queue.length > 0) {
        const ch = queue.shift();
        result.push(ch);
        for (const neighbor of adj.get(ch)) {
            inDegree.set(neighbor, inDegree.get(neighbor) - 1);
            if (inDegree.get(neighbor) === 0) queue.push(neighbor);
        }
    }

    if (result.length !== inDegree.size) return "";
    return result.join("");
}

module.exports = { alienOrder };
