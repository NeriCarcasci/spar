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

fn validate(node: &Option<Rc<RefCell<TreeNode>>>, lo: i64, hi: i64) -> bool {
    match node {
        None => true,
        Some(n) => {
            let n = n.borrow();
            let val = n.val as i64;
            if val <= lo || val >= hi {
                return false;
            }
            validate(&n.left, lo, val) && validate(&n.right, val, hi)
        }
    }
}

pub fn is_valid_bst(root: Option<Rc<RefCell<TreeNode>>>) -> bool {
    validate(&root, i64::MIN, i64::MAX)
}
