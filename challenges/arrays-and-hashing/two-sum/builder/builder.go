package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type testCase struct {
	Input    []int
	Target   int
	Expected []int
}

func main() {
	testsPath := "../tests.yaml"
	if len(os.Args) > 2 {
		testsPath = os.Args[2]
	}

	tests, err := parseTests(testsPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	failed := false
	for i, tc := range tests {
		got := twoSum(tc.Input, tc.Target)
		if equalSlices(got, tc.Expected) {
			fmt.Printf("PASS %d\n", i+1)
			continue
		}
		failed = true
		fmt.Printf("FAIL %d got=%s expected=%s\n", i+1, formatSlice(got), formatSlice(tc.Expected))
	}

	if failed {
		os.Exit(1)
	}
}

func parseTests(path string) ([]testCase, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tests []testCase
	var current *testCase
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || trimmed == "visible:" || trimmed == "hidden:" {
			continue
		}
		if strings.HasPrefix(trimmed, "- ") {
			if current != nil {
				tests = append(tests, *current)
			}
			current = &testCase{}
			trimmed = strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))
		}
		if current == nil || !strings.Contains(trimmed, ":") {
			continue
		}
		key, value := splitKV(trimmed)
		switch key {
		case "input":
			vals, err := parseIntSlice(value)
			if err != nil {
				return nil, err
			}
			current.Input = vals
		case "target":
			n, err := strconv.Atoi(strings.TrimSpace(value))
			if err != nil {
				return nil, err
			}
			current.Target = n
		case "expected":
			vals, err := parseIntSlice(value)
			if err != nil {
				return nil, err
			}
			current.Expected = vals
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if current != nil {
		tests = append(tests, *current)
	}
	return tests, nil
}

func splitKV(line string) (string, string) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) == 1 {
		return strings.TrimSpace(parts[0]), ""
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
}

func parseIntSlice(value string) ([]int, error) {
	raw := strings.TrimSpace(value)
	raw = strings.TrimPrefix(raw, "[")
	raw = strings.TrimSuffix(raw, "]")
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []int{}, nil
	}
	parts := strings.Split(raw, ",")
	result := make([]int, 0, len(parts))
	for _, part := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return nil, err
		}
		result = append(result, n)
	}
	return result, nil
}

func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func formatSlice(values []int) string {
	if len(values) == 0 {
		return "[]"
	}
	parts := make([]string, 0, len(values))
	for _, v := range values {
		parts = append(parts, strconv.Itoa(v))
	}
	return "[" + strings.Join(parts, ",") + "]"
}
