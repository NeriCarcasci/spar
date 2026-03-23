package solution

func LongestConsecutive(nums []int) int {
	numSet := make(map[int]struct{}, len(nums))
	for _, n := range nums {
		numSet[n] = struct{}{}
	}
	longest := 0
	for n := range numSet {
		if _, found := numSet[n-1]; found {
			continue
		}
		length := 1
		for {
			if _, found := numSet[n+length]; !found {
				break
			}
			length++
		}
		if length > longest {
			longest = length
		}
	}
	return longest
}
