function evalRPN(tokens) {
    const stack = [];
    for (const token of tokens) {
        if (["+", "-", "*", "/"].includes(token)) {
            const b = stack.pop();
            const a = stack.pop();
            switch (token) {
                case "+": stack.push(a + b); break;
                case "-": stack.push(a - b); break;
                case "*": stack.push(a * b); break;
                case "/": stack.push(Math.trunc(a / b)); break;
            }
        } else {
            stack.push(parseInt(token, 10));
        }
    }
    return stack[0];
}
module.exports = { evalRPN };
