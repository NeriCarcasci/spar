#[derive(Debug, PartialEq)]
pub struct ListNode {
    pub val: i32,
    pub next: Option<Box<ListNode>>,
}

impl ListNode {
    pub fn new(val: i32) -> Self {
        ListNode { val, next: None }
    }
}

pub fn reorder_list(head: &mut Option<Box<ListNode>>) {
    let mut nodes = Vec::new();
    let mut curr = head.take();
    while let Some(mut node) = curr {
        curr = node.next.take();
        nodes.push(node);
    }
    let mut left = 0;
    let mut right = nodes.len().saturating_sub(1);
    let mut dummy = Box::new(ListNode::new(0));
    let mut tail = &mut dummy;
    while left <= right && right < nodes.len() {
        let node = std::mem::replace(&mut nodes[left], Box::new(ListNode::new(0)));
        tail.next = Some(node);
        tail = tail.next.as_mut().unwrap();
        if left != right {
            let node = std::mem::replace(&mut nodes[right], Box::new(ListNode::new(0)));
            tail.next = Some(node);
            tail = tail.next.as_mut().unwrap();
        }
        left += 1;
        if right == 0 { break; }
        right -= 1;
    }
    *head = dummy.next;
}
