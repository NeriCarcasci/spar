pub fn search_rotated(nums: Vec<i32>, target: i32) -> i32 {
    let (mut left, mut right) = (0i32, nums.len() as i32 - 1);
    while left <= right {
        let mid = left + (right - left) / 2;
        let m = mid as usize;
        if nums[m] == target { return mid; }
        if nums[left as usize] <= nums[m] {
            if nums[left as usize] <= target && target < nums[m] { right = mid - 1; }
            else { left = mid + 1; }
        } else {
            if nums[m] < target && target <= nums[right as usize] { left = mid + 1; }
            else { right = mid - 1; }
        }
    }
    -1
}
