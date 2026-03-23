use std::cell::RefCell;
use std::collections::VecDeque;
use std::rc::Rc;

#[derive(Debug, PartialEq)]
pub struct TreeNode {
    pub val: i32,
    pub left: Option<Rc<RefCell<TreeNode>>>,
    pub right: Option<Rc<RefCell<TreeNode>>>,
}

impl TreeNode {
    pub fn new(val: i32) -> Self {
        TreeNode { val, left: None, right: None }
    }
}

pub fn level_order(root: Option<Rc<RefCell<TreeNode>>>) -> Vec<Vec<i32>> {
    let mut result = vec![];
    if root.is_none() {
        return result;
    }
    let mut queue = VecDeque::new();
    queue.push_back(root.unwrap());
    while !queue.is_empty() {
        let size = queue.len();
        let mut level = vec![];
        for _ in 0..size {
            let node = queue.pop_front().unwrap();
            let node = node.borrow();
            level.push(node.val);
            if let Some(ref left) = node.left {
                queue.push_back(Rc::clone(left));
            }
            if let Some(ref right) = node.right {
                queue.push_back(Rc::clone(right));
            }
        }
        result.push(level);
    }
    result
}
