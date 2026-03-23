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

fn is_same(p: &Option<Rc<RefCell<TreeNode>>>, q: &Option<Rc<RefCell<TreeNode>>>) -> bool {
    match (p, q) {
        (None, None) => true,
        (Some(a), Some(b)) => {
            let a = a.borrow();
            let b = b.borrow();
            a.val == b.val && is_same(&a.left, &b.left) && is_same(&a.right, &b.right)
        }
        _ => false,
    }
}

pub fn is_subtree(root: Option<Rc<RefCell<TreeNode>>>, sub_root: Option<Rc<RefCell<TreeNode>>>) -> bool {
    match &root {
        None => false,
        Some(node) => {
            if is_same(&root, &sub_root) {
                return true;
            }
            let node = node.borrow();
            is_subtree(node.left.clone(), sub_root.clone()) || is_subtree(node.right.clone(), sub_root)
        }
    }
}
