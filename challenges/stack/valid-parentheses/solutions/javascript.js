function isValid(s) {
    const stack = [];
    const pairs = { ')': '(', '}': '{', ']': '[' };
    for (const ch of s) {
        if (pairs[ch]) {
            if (!stack.length || stack[stack.length - 1] !== pairs[ch]) return false;
            stack.pop();
        } else {
            stack.push(ch);
        }
    }
    return stack.length === 0;
}
module.exports = { isValid };
