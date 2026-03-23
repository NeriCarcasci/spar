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

pub fn serialize(root: Option<Rc<RefCell<TreeNode>>>) -> String {
    String::new()
}

pub fn deserialize(data: String) -> Option<Rc<RefCell<TreeNode>>> {
    None
}
