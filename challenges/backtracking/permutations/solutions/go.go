package solution

import "sort"

func Permute(nums []int) [][]int {
	sort.Ints(nums)
	var result [][]int
	used := make([]bool, len(nums))
	var backtrack func(path []int)
	backtrack = func(path []int) {
		if len(path) == len(nums) {
			perm := make([]int, len(path))
			copy(perm, path)
			result = append(result, perm)
			return
		}
		for i := 0; i < len(nums); i++ {
			if used[i] {
				continue
			}
			used[i] = true
			backtrack(append(path, nums[i]))
			used[i] = false
		}
	}
	backtrack(nil)
	return result
}
