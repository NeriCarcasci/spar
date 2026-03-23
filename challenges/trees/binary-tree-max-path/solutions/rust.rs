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

fn dfs(node: &Option<Rc<RefCell<TreeNode>>>, result: &mut i32) -> i32 {
    match node {
        None => 0,
        Some(n) => {
            let n = n.borrow();
            let left_gain = dfs(&n.left, result).max(0);
            let right_gain = dfs(&n.right, result).max(0);
            *result = (*result).max(n.val + left_gain + right_gain);
            n.val + left_gain.max(right_gain)
        }
    }
}

pub fn max_path_sum(root: Option<Rc<RefCell<TreeNode>>>) -> i32 {
    let mut result = i32::MIN;
    dfs(&root, &mut result);
    result
}
