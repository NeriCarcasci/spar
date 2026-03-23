package solution

func LeastInterval(tasks []byte, n int) int {
	freq := [26]int{}
	for _, t := range tasks {
		freq[t-'A']++
	}
	maxFreq := 0
	for _, f := range freq {
		if f > maxFreq {
			maxFreq = f
		}
	}
	countMax := 0
	for _, f := range freq {
		if f == maxFreq {
			countMax++
		}
	}
	result := (maxFreq-1)*(n+1) + countMax
	if len(tasks) > result {
		return len(tasks)
	}
	return result
}
