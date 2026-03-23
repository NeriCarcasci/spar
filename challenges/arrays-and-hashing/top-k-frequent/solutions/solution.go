package solution

import "sort"

func TopKFrequent(nums []int, k int) []int {
	freq := make(map[int]int)
	for _, n := range nums {
		freq[n]++
	}

	buckets := make([][]int, len(nums)+1)
	for num, count := range freq {
		buckets[count] = append(buckets[count], num)
	}

	result := make([]int, 0, k)
	for i := len(buckets) - 1; i > 0 && len(result) < k; i-- {
		result = append(result, buckets[i]...)
	}
	result = result[:k]
	sort.Ints(result)
	return result
}
