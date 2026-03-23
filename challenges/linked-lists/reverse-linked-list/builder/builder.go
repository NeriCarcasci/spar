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
	Expected []int
}

func main() {
	testsPath := "../tests.yaml"
	challengePath := "../challenge.yaml"
	if len(os.Args) > 2 {
		testsPath = os.Args[2]
	}
	if len(os.Args) > 3 {
		challengePath = os.Args[3]
	}

	tests, err := parseTests(testsPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	inputType, outputType, err := parseTypes(challengePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	failed := false
	for i, tc := range tests {
		var got []int
		if inputType == "linked-list" {
			head := sliceToList(tc.Input)
			out := reverseList(head)
			if outputType == "linked-list" {
				got = listToSlice(out)
			}
		}
		if outputType != "linked-list" {
			got = tc.Input
		}

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

func parseTypes(path string) (string, string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "array", "array", err
	}
	defer file.Close()

	inputType := "array"
	outputType := "array"
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, " ") {
			continue
		}
		trimmed := strings.TrimSpace(line)
		if !strings.Contains(trimmed, ":") {
			continue
		}
		key, value := splitKV(trimmed)
		value = strings.Trim(value, "\"'")
		if key == "input_type" {
			inputType = value
		}
		if key == "output_type" {
			outputType = value
		}
	}
	if err := scanner.Err(); err != nil {
		return "array", "array", err
	}
	return inputType, outputType, nil
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

func sliceToList(values []int) *ListNode {
	var head *ListNode
	var tail *ListNode
	for _, v := range values {
		n := &ListNode{Val: v}
		if head == nil {
			head = n
			tail = n
			continue
		}
		tail.Next = n
		tail = n
	}
	return head
}

func listToSlice(head *ListNode) []int {
	result := []int{}
	for node := head; node != nil; node = node.Next {
		result = append(result, node.Val)
	}
	return result
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
