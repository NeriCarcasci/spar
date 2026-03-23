class ListNode {
    constructor(val = 0, next = null) {
        this.val = val;
        this.next = next;
    }
}

function reorderList(head) {
    if (!head || !head.next) return;
    let slow = head, fast = head;
    while (fast.next && fast.next.next) {
        slow = slow.next;
        fast = fast.next.next;
    }
    let second = slow.next;
    slow.next = null;
    let prev = null;
    while (second) {
        const next = second.next;
        second.next = prev;
        prev = second;
        second = next;
    }
    let first = head;
    second = prev;
    while (second) {
        const [tmp1, tmp2] = [first.next, second.next];
        first.next = second;
        second.next = tmp1;
        first = tmp1;
        second = tmp2;
    }
}

module.exports = { ListNode, reorderList };
