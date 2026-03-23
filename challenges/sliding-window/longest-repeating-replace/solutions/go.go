package solution

func CharacterReplacement(s string, k int) int {
	var counts [26]int
	left, maxFreq, longest := 0, 0, 0
	for right := 0; right < len(s); right++ {
		counts[s[right]-'A']++
		if counts[s[right]-'A'] > maxFreq {
			maxFreq = counts[s[right]-'A']
		}
		for (right-left+1)-maxFreq > k {
			counts[s[left]-'A']--
			left++
		}
		if length := right - left + 1; length > longest {
			longest = length
		}
	}
	return longest
}
