function minWindow(s, t) {
    const need = new Map();
    for (const ch of t) need.set(ch, (need.get(ch) || 0) + 1);
    let missing = need.size;
    let left = 0, bestLeft = 0, bestRight = s.length, found = false;
    for (let right = 0; right < s.length; right++) {
        need.set(s[right], (need.get(s[right]) || 0) - 1);
        if (need.get(s[right]) === 0) missing--;
        while (missing === 0) {
            if (right - left < bestRight - bestLeft) {
                bestLeft = left;
                bestRight = right;
                found = true;
            }
            need.set(s[left], need.get(s[left]) + 1);
            if (need.get(s[left]) > 0) missing++;
            left++;
        }
    }
    return found ? s.substring(bestLeft, bestRight + 1) : "";
}
module.exports = { minWindow };
