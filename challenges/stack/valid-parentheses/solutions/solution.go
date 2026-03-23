package solution

func IsValid(s string) bool {
	stack := []byte{}
	pairs := map[byte]byte{')': '(', '}': '{', ']': '['}
	for i := 0; i < len(s); i++ {
		if match, ok := pairs[s[i]]; ok {
			if len(stack) == 0 || stack[len(stack)-1] != match {
				return false
			}
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, s[i])
		}
	}
	return len(stack) == 0
}
