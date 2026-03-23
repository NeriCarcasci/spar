use std::collections::HashSet;

#[derive(Debug)]
pub struct ListNode {
    pub val: i32,
    pub next: Option<Box<ListNode>>,
}

pub fn has_cycle(head: &Option<Box<ListNode>>) -> bool {
    let mut seen = HashSet::new();
    let mut current = head;
    while let Some(node) = current {
        let ptr = &**node as *const ListNode as usize;
        if !seen.insert(ptr) {
            return true;
        }
        current = &node.next;
    }
    false
}
