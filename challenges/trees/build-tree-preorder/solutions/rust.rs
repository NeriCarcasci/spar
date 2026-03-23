use std::cell::RefCell;
use std::collections::HashMap;
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

fn helper(
    preorder: &[i32],
    pre_idx: &mut usize,
    in_left: i32,
    in_right: i32,
    in_map: &HashMap<i32, i32>,
) -> Option<Rc<RefCell<TreeNode>>> {
    if in_left > in_right {
        return None;
    }
    let root_val = preorder[*pre_idx];
    *pre_idx += 1;
    let mid = in_map[&root_val];
    let left = helper(preorder, pre_idx, in_left, mid - 1, in_map);
    let right = helper(preorder, pre_idx, mid + 1, in_right, in_map);
    Some(Rc::new(RefCell::new(TreeNode {
        val: root_val,
        left,
        right,
    })))
}

pub fn build_tree(preorder: Vec<i32>, inorder: Vec<i32>) -> Option<Rc<RefCell<TreeNode>>> {
    if preorder.is_empty() {
        return None;
    }
    let mut in_map = HashMap::new();
    for (i, &v) in inorder.iter().enumerate() {
        in_map.insert(v, i as i32);
    }
    let mut pre_idx = 0;
    helper(&preorder, &mut pre_idx, 0, inorder.len() as i32 - 1, &in_map)
}
