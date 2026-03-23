package solution

import "sort"

func CombinationSum(candidates []int, target int) [][]int {
	sort.Ints(candidates)
	var result [][]int
	var backtrack func(start, remaining int, path []int)
	backtrack = func(start, remaining int, path []int) {
		if remaining == 0 {
			combo := make([]int, len(path))
			copy(combo, path)
			result = append(result, combo)
			return
		}
		for i := start; i < len(candidates); i++ {
			if candidates[i] > remaining {
				break
			}
			backtrack(i, remaining-candidates[i], append(path, candidates[i]))
		}
	}
	backtrack(0, target, nil)
	return result
}
