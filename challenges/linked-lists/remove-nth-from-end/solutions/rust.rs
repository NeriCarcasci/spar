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

pub fn remove_nth_from_end(head: Option<Box<ListNode>>, n: i32) -> Option<Box<ListNode>> {
    let mut dummy = Box::new(ListNode { val: 0, next: head });
    let mut length = 0;
    let mut curr = &dummy.next;
    while let Some(node) = curr {
        length += 1;
        curr = &node.next;
    }
    let target = length - n as usize;
    let mut curr = &mut dummy;
    for _ in 0..target {
        curr = curr.next.as_mut().unwrap();
    }
    let next = curr.next.as_mut().unwrap().next.take();
    curr.next = next;
    dummy.next
}
