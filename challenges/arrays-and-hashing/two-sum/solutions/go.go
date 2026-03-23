package main

func twoSum(nums []int, target int) []int {
	seen := make(map[int]int)
	for i, num := range nums {
		want := target - num
		if j, ok := seen[want]; ok {
			return []int{j, i}
		}
		seen[num] = i
	}
	return []int{}
}
