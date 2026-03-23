package solution

func MaxSlidingWindow(nums []int, k int) []int {
	var dq []int
	var result []int
	for i, n := range nums {
		for len(dq) > 0 && dq[0] < i-k+1 {
			dq = dq[1:]
		}
		for len(dq) > 0 && nums[dq[len(dq)-1]] <= n {
			dq = dq[:len(dq)-1]
		}
		dq = append(dq, i)
		if i >= k-1 {
			result = append(result, nums[dq[0]])
		}
	}
	return result
}
