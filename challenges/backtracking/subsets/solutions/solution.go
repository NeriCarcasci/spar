package solution

import "sort"

func Subsets(nums []int) [][]int {
	sort.Ints(nums)
	var result [][]int
	var backtrack func(start int, path []int)
	backtrack = func(start int, path []int) {
		sub := make([]int, len(path))
		copy(sub, path)
		result = append(result, sub)
		for i := start; i < len(nums); i++ {
			backtrack(i+1, append(path, nums[i]))
		}
	}
	backtrack(0, nil)
	return result
}
