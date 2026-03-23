package solution

func AlienOrder(words []string) string {
	adj := map[byte]map[byte]bool{}
	inDegree := map[byte]int{}

	for _, word := range words {
		for i := 0; i < len(word); i++ {
			if _, ok := inDegree[word[i]]; !ok {
				inDegree[word[i]] = 0
				adj[word[i]] = map[byte]bool{}
			}
		}
	}

	for i := 0; i < len(words)-1; i++ {
		w1, w2 := words[i], words[i+1]
		minLen := len(w1)
		if len(w2) < minLen {
			minLen = len(w2)
		}
		if len(w1) > len(w2) && w1[:minLen] == w2[:minLen] {
			return ""
		}
		for j := 0; j < minLen; j++ {
			if w1[j] != w2[j] {
				if !adj[w1[j]][w2[j]] {
					adj[w1[j]][w2[j]] = true
					inDegree[w2[j]]++
				}
				break
			}
		}
	}

	queue := []byte{}
	for ch, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, ch)
		}
	}

	result := []byte{}
	for len(queue) > 0 {
		ch := queue[0]
		queue = queue[1:]
		result = append(result, ch)
		for neighbor := range adj[ch] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	if len(result) != len(inDegree) {
		return ""
	}
	return string(result)
}
