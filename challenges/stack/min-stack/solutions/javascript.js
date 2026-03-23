class MinStack {
    constructor() {
        this.stack = [];
    }

    push(val) {
        const currentMin = this.stack.length ? Math.min(val, this.stack[this.stack.length - 1][1]) : val;
        this.stack.push([val, currentMin]);
    }

    pop() {
        this.stack.pop();
    }

    top() {
        return this.stack[this.stack.length - 1][0];
    }

    getMin() {
        return this.stack[this.stack.length - 1][1];
    }
}

module.exports = { MinStack };
