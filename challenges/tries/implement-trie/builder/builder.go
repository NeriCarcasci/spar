package main

import (
    "encoding/json"
    "fmt"
    "math"
    "os"
    "reflect"
    "sort"
)

const testsJSON = "[{\"input\":{\"operations\":[\"insert\",\"search\",\"search\",\"starts_with\",\"insert\",\"search\"],\"values\":[\"apple\",\"apple\",\"app\",\"app\",\"app\",\"app\"]},\"expected\":[null,true,false,true,null,true],\"visible\":true},{\"input\":{\"operations\":[\"insert\",\"search\",\"starts_with\"],\"values\":[\"hello\",\"hello\",\"hell\"]},\"expected\":[null,true,true],\"visible\":true},{\"input\":{\"operations\":[\"search\",\"starts_with\"],\"values\":[\"empty\",\"empty\"]},\"expected\":[false,false],\"visible\":false},{\"input\":{\"operations\":[\"insert\",\"insert\",\"search\",\"search\"],\"values\":[\"a\",\"b\",\"a\",\"b\"]},\"expected\":[null,null,true,true],\"visible\":false},{\"input\":{\"operations\":[\"insert\",\"insert\",\"insert\",\"starts_with\",\"starts_with\",\"search\"],\"values\":[\"abc\",\"ab\",\"a\",\"a\",\"abc\",\"ab\"]},\"expected\":[null,null,null,true,true,true],\"visible\":false},{\"input\":{\"operations\":[\"insert\",\"search\",\"starts_with\",\"search\"],\"values\":[\"prefix\",\"prefix\",\"pre\",\"pre\"]},\"expected\":[null,true,true,false],\"visible\":false},{\"input\":{\"operations\":[\"insert\",\"insert\",\"search\",\"search\",\"starts_with\"],\"values\":[\"app\",\"apple\",\"app\",\"apple\",\"appl\"]},\"expected\":[null,null,true,true,true],\"visible\":false}]"
const compareMode = "exact"

func main() {
    testsPath := "../tests.yaml"
    challengePath := "../challenge.yaml"
    if len(os.Args) > 3 {
        testsPath = os.Args[3]
    }
    if len(os.Args) > 4 {
        challengePath = os.Args[4]
    }
    _, _ = os.ReadFile(testsPath)
    _, _ = os.ReadFile(challengePath)

    var tests []map[string]any
    if err := json.Unmarshal([]byte(testsJSON), &tests); err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        os.Exit(1)
    }

    argKeys := []string{"operations", "values"}
    _ = argKeys
    failed := false
    for i, test := range tests {
        var gotValue any

        inp := getInputMap(test)
        obj := NewTrie()
        operations := toStringSlice(inp["operations"])
        values := toAnySlice(inp["values"])
        out := make([]any, 0, len(operations))
        for i, op := range operations {
            val := ""
            if i < len(values) {
                val = toString(values[i])
            }
            switch op {
            case "insert":
                obj.Insert(val)
                out = append(out, nil)
            case "search":
                out = append(out, obj.Search(val))
            case "starts_with":
                out = append(out, obj.StartsWith(val))
            default:
                panic("unknown op")
            }
        }
        gotValue = normalizeValue(out)

        expected := test["expected"]
        expectedLength, hasExpectedLength := test["expected_length"]
        passed := false
        if hasExpectedLength {
            if arr, ok := gotValue.([]any); ok {
                passed = len(arr) == toInt(expectedLength)
            } else {
                passed = false
            }
        } else if compareMode == "pair_target_1idx" {
            inputObj := getInputMap(test)
            numbers := toAnySlice(inputObj["numbers"])
            target := toInt(inputObj["target"])
            if pair, ok := gotValue.([]any); ok && len(pair) == 2 {
                i1 := toInt(pair[0]) - 1
                i2 := toInt(pair[1]) - 1
                if i1 >= 0 && i2 >= 0 && i1 < len(numbers) && i2 < len(numbers) && i1 != i2 {
                    passed = toInt(numbers[i1])+toInt(numbers[i2]) == target
                } else {
                    passed = false
                }
            } else {
                passed = false
            }
        } else {
            passed = equalValues(gotValue, expected, compareMode)
        }
        idx := i + 1
        if passed {
            fmt.Printf("PASS %d\n", idx)
        } else {
            failed = true
            if hasExpectedLength {
                fmt.Printf("FAIL %d got=%s expected=%s\n", idx, render(gotValue), render(normalizeValue(expectedLength)))
            } else {
                fmt.Printf("FAIL %d got=%s expected=%s\n", idx, render(gotValue), render(expected))
            }
        }
    }
    if failed {
        os.Exit(1)
    }
}

func getInputMap(test map[string]any) map[string]any {
    if input, ok := test["input"]; ok {
        if m, ok := input.(map[string]any); ok {
            return m
        }
    }
    return map[string]any{}
}

func getArg(test map[string]any, argKeys []string, idx int) any {
    input, hasInput := test["input"]
    if hasInput {
        if m, ok := input.(map[string]any); ok {
            if idx < len(argKeys) {
                return m[argKeys[idx]]
            }
            return nil
        }
        if arr, ok := input.([]any); ok {
            if len(argKeys) > 1 {
                if idx == 0 {
                    return arr
                }
                key := argKeys[idx]
                if v, ok := test[key]; ok {
                    return v
                }
                if idx < len(arr) {
                    return arr[idx]
                }
                return nil
            }
            return arr
        }
        return input
    }
    if idx < len(argKeys) {
        return test[argKeys[idx]]
    }
    return nil
}

func toAnySlice(v any) []any {
    if v == nil {
        return []any{}
    }
    if arr, ok := v.([]any); ok {
        return arr
    }
    return []any{}
}

func toInt(v any) int {
    switch value := v.(type) {
    case nil:
        return 0
    case int:
        return value
    case int64:
        return int(value)
    case float64:
        return int(value)
    case json.Number:
        n, _ := value.Int64()
        return int(n)
    default:
        return 0
    }
}

func toFloat(v any) float64 {
    switch value := v.(type) {
    case nil:
        return 0
    case float64:
        return value
    case int:
        return float64(value)
    case json.Number:
        n, _ := value.Float64()
        return n
    default:
        return 0
    }
}

func toString(v any) string {
    switch value := v.(type) {
    case nil:
        return ""
    case string:
        return value
    default:
        return fmt.Sprint(value)
    }
}

func toBool(v any) bool {
    if b, ok := v.(bool); ok {
        return b
    }
    return false
}

func toIntSlice(v any) []int {
    raw := toAnySlice(v)
    out := make([]int, 0, len(raw))
    for _, item := range raw {
        out = append(out, toInt(item))
    }
    return out
}

func toIntMatrix(v any) [][]int {
    raw := toAnySlice(v)
    out := make([][]int, 0, len(raw))
    for _, item := range raw {
        out = append(out, toIntSlice(item))
    }
    return out
}

func toStringSlice(v any) []string {
    raw := toAnySlice(v)
    out := make([]string, 0, len(raw))
    for _, item := range raw {
        out = append(out, toString(item))
    }
    return out
}

func toStringMatrix(v any) [][]string {
    raw := toAnySlice(v)
    out := make([][]string, 0, len(raw))
    for _, item := range raw {
        out = append(out, toStringSlice(item))
    }
    return out
}

func toByteSlice(v any) []byte {
    raw := toAnySlice(v)
    out := make([]byte, 0, len(raw))
    for _, item := range raw {
        s := toString(item)
        if s == "" {
            out = append(out, 0)
        } else {
            out = append(out, s[0])
        }
    }
    return out
}

func toByteMatrix(v any) [][]byte {
    raw := toAnySlice(v)
    out := make([][]byte, 0, len(raw))
    for _, item := range raw {
        out = append(out, toByteSlice(item))
    }
    return out
}

func normalizeValue(v any) any {
    data, err := json.Marshal(v)
    if err != nil {
        return v
    }
    var out any
    if err := json.Unmarshal(data, &out); err != nil {
        return v
    }
    return out
}

func render(v any) string {
    normalized := normalizeValue(v)
    data, err := json.Marshal(normalized)
    if err != nil {
        return "null"
    }
    return string(data)
}

func sortPrimitiveSlice(values []any) []any {
    out := append([]any{}, values...)
    sort.Slice(out, func(i, j int) bool {
        return fmt.Sprint(out[i]) < fmt.Sprint(out[j])
    })
    return out
}

func canonical(v any, mode string) any {
    if mode == "pair_unordered" || mode == "list_unordered" || mode == "strings_unordered" {
        if arr, ok := v.([]any); ok {
            return sortPrimitiveSlice(arr)
        }
        return v
    }
    if mode == "groups_unordered" || mode == "nested_unordered" {
        arr, ok := v.([]any)
        if !ok {
            return v
        }
        outer := make([]any, 0, len(arr))
        for _, item := range arr {
            if inner, ok := item.([]any); ok {
                outer = append(outer, sortPrimitiveSlice(inner))
            } else {
                outer = append(outer, item)
            }
        }
        sort.Slice(outer, func(i, j int) bool {
            return render(outer[i]) < render(outer[j])
        })
        return outer
    }
    return v
}

func equalFloatSequence(got, expected any) bool {
    ga, ok1 := got.([]any)
    ea, ok2 := expected.([]any)
    if !ok1 || !ok2 || len(ga) != len(ea) {
        return false
    }
    for i := range ga {
        if ga[i] == nil && ea[i] == nil {
            continue
        }
        if ga[i] == nil || ea[i] == nil {
            return false
        }
        if math.Abs(toFloat(ga[i])-toFloat(ea[i])) > 1e-9 {
            return false
        }
    }
    return true
}

func equalValues(got, expected any, mode string) bool {
    g := normalizeValue(got)
    e := normalizeValue(expected)
    if mode == "pair_target_1idx" {
        return false
    }
    if mode == "float_sequence" {
        return equalFloatSequence(g, e)
    }
    cg := canonical(g, mode)
    ce := canonical(e, mode)
    if reflect.DeepEqual(cg, ce) {
        return true
    }
    if cg == nil {
        if arr, ok := ce.([]any); ok && len(arr) == 0 {
            return true
        }
    }
    if ce == nil {
        if arr, ok := cg.([]any); ok && len(arr) == 0 {
            return true
        }
    }
    return false
}




