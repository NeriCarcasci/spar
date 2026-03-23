package solution

func MinWindow(s string, t string) string {
	need := make(map[byte]int)
	for i := range t {
		need[t[i]]++
	}
	missing := len(need)
	left := 0
	bestLeft, bestRight := 0, len(s)
	found := false
	for right := 0; right < len(s); right++ {
		need[s[right]]--
		if need[s[right]] == 0 {
			missing--
		}
		for missing == 0 {
			if right-left < bestRight-bestLeft {
				bestLeft, bestRight = left, right
				found = true
			}
			need[s[left]]++
			if need[s[left]] > 0 {
				missing++
			}
			left++
		}
	}
	if !found {
		return ""
	}
	return s[bestLeft : bestRight+1]
}
