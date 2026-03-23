class Heap {
    constructor(comparator) {
        this.data = [];
        this.cmp = comparator;
    }
    size() { return this.data.length; }
    peek() { return this.data[0]; }
    push(val) {
        this.data.push(val);
        this._bubbleUp(this.data.length - 1);
    }
    pop() {
        const top = this.data[0];
        const last = this.data.pop();
        if (this.data.length > 0) {
            this.data[0] = last;
            this._sinkDown(0);
        }
        return top;
    }
    _bubbleUp(i) {
        while (i > 0) {
            const p = (i - 1) >> 1;
            if (this.cmp(this.data[i], this.data[p]) >= 0) break;
            [this.data[i], this.data[p]] = [this.data[p], this.data[i]];
            i = p;
        }
    }
    _sinkDown(i) {
        const n = this.data.length;
        while (true) {
            let s = i;
            const l = 2 * i + 1, r = 2 * i + 2;
            if (l < n && this.cmp(this.data[l], this.data[s]) < 0) s = l;
            if (r < n && this.cmp(this.data[r], this.data[s]) < 0) s = r;
            if (s === i) break;
            [this.data[i], this.data[s]] = [this.data[s], this.data[i]];
            i = s;
        }
    }
}

class MedianFinder {
    constructor() {
        this.lo = new Heap((a, b) => b - a);
        this.hi = new Heap((a, b) => a - b);
    }

    addNum(num) {
        this.lo.push(num);
        this.hi.push(this.lo.pop());
        if (this.hi.size() > this.lo.size()) {
            this.lo.push(this.hi.pop());
        }
    }

    findMedian() {
        if (this.lo.size() > this.hi.size()) {
            return this.lo.peek();
        }
        return (this.lo.peek() + this.hi.peek()) / 2;
    }
}

module.exports = { MedianFinder };
