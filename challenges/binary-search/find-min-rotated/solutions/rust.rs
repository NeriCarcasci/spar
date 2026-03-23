pub fn find_min(nums: Vec<i32>) -> i32 {
    let (mut left, mut right) = (0usize, nums.len() - 1);
    while left < right {
        let mid = left + (right - left) / 2;
        if nums[mid] > nums[right] { left = mid + 1; }
        else { right = mid; }
    }
    nums[left]
}
