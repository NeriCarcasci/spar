package solution

func LengthOfLongestSubstring(s string) int {
	lastSeen := make(map[byte]int)
	start := 0
	longest := 0
	for i := 0; i < len(s); i++ {
		if j, ok := lastSeen[s[i]]; ok && j >= start {
			start = j + 1
		}
		lastSeen[s[i]] = i
		if length := i - start + 1; length > longest {
			longest = length
		}
	}
	return longest
}
