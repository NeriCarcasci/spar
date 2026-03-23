pub fn largest_rectangle_area(heights: Vec<i32>) -> i32 {
    let mut stack: Vec<(usize, i32)> = Vec::new();
    let mut max_area = 0;
    for (i, &h) in heights.iter().enumerate() {
        let mut start = i;
        while stack.last().map_or(false, |&(_, sh)| sh > h) {
            let (idx, sh) = stack.pop().unwrap();
            max_area = max_area.max(sh * (i - idx) as i32);
            start = idx;
        }
        stack.push((start, h));
    }
    for (idx, h) in stack {
        max_area = max_area.max(h * (heights.len() - idx) as i32);
    }
    max_area
}
