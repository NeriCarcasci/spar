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
    fn dfs(node: &Option<Rc<RefCell<TreeNode>>>, result: &mut Vec<String>) {
        match node {
            None => result.push("N".to_string()),
            Some(n) => {
                let n = n.borrow();
                result.push(n.val.to_string());
                dfs(&n.left, result);
                dfs(&n.right, result);
            }
        }
    }
    let mut result = vec![];
    dfs(&root, &mut result);
    result.join(",")
}

pub fn deserialize(data: String) -> Option<Rc<RefCell<TreeNode>>> {
    fn dfs(vals: &[&str], idx: &mut usize) -> Option<Rc<RefCell<TreeNode>>> {
        if *idx >= vals.len() || vals[*idx] == "N" {
            *idx += 1;
            return None;
        }
        let val: i32 = vals[*idx].parse().unwrap();
        *idx += 1;
        let left = dfs(vals, idx);
        let right = dfs(vals, idx);
        Some(Rc::new(RefCell::new(TreeNode { val, left, right })))
    }
    let vals: Vec<&str> = data.split(',').collect();
    let mut idx = 0;
    dfs(&vals, &mut idx)
}
