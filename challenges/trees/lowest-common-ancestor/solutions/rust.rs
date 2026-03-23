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

pub fn lowest_common_ancestor(root: Option<Rc<RefCell<TreeNode>>>, p: i32, q: i32) -> i32 {
    let mut cur = root;
    while let Some(node) = cur {
        let val = node.borrow().val;
        if p < val && q < val {
            cur = node.borrow().left.clone();
        } else if p > val && q > val {
            cur = node.borrow().right.clone();
        } else {
            return val;
        }
    }
    -1
}
