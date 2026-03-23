use std::cell::RefCell;
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

pub fn kth_smallest(root: Option<Rc<RefCell<TreeNode>>>, k: i32) -> i32 {
    let mut stack: Vec<Rc<RefCell<TreeNode>>> = vec![];
    let mut cur = root;
    let mut count = 0;
    loop {
        while let Some(node) = cur {
            cur = node.borrow().left.clone();
            stack.push(node);
        }
        if let Some(node) = stack.pop() {
            count += 1;
            if count == k {
                return node.borrow().val;
            }
            cur = node.borrow().right.clone();
        } else {
            break;
        }
    }
    -1
}
