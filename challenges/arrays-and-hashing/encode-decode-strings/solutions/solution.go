package solution

import (
	"strconv"
	"strings"
)

func Encode(strs []string) string {
	var b strings.Builder
	for _, s := range strs {
		b.WriteString(strconv.Itoa(len(s)))
		b.WriteByte('#')
		b.WriteString(s)
	}
	return b.String()
}

func Decode(s string) []string {
	var result []string
	i := 0
	for i < len(s) {
		j := i
		for s[j] != '#' {
			j++
		}
		length, _ := strconv.Atoi(s[i:j])
		result = append(result, s[j+1:j+1+length])
		i = j + 1 + length
	}
	return result
}
