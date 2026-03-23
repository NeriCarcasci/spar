package runner

import (
	"strings"
	"testing"
)

func TestParseResults(t *testing.T) {
	stdout := strings.Join([]string{
		"PASS 1",
		"FAIL 2 got=[1,0] expected=[0,1]",
		"PASS 3",
	}, "\n")

	results := ParseResults(stdout)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if !results[0].Passed || results[0].Index != 1 {
		t.Fatalf("unexpected first result: %#v", results[0])
	}
	if results[1].Passed || results[1].Got != "[1,0]" || results[1].Expected != "[0,1]" {
		t.Fatalf("unexpected second result: %#v", results[1])
	}
}

func TestNormalizeLanguage(t *testing.T) {
	cases := map[string]string{
		"python":     "python",
		"JS":         "javascript",
		"javascript": "javascript",
		"c++":        "cpp",
		"Cpp":        "cpp",
	}
	for input, expected := range cases {
		if got := normalizeLanguage(input); got != expected {
			t.Fatalf("normalizeLanguage(%q) = %q, expected %q", input, got, expected)
		}
	}
}

func TestSpecForLanguage(t *testing.T) {
	spec, err := specForLanguage("go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if spec.BuilderFile != "builder.go" || spec.SolutionFile != "solution.go" {
		t.Fatalf("unexpected spec: %#v", spec)
	}
	if _, err := specForLanguage("unknown"); err == nil {
		t.Fatal("expected error for unknown language")
	}
}

func TestClassifyFailure(t *testing.T) {
	compileErr, runtimeErr := classifyFailure("go", nil, "", "syntax error", "", "")
	if compileErr == "" || runtimeErr != "" {
		t.Fatalf("unexpected classification for compile error: compile=%q runtime=%q", compileErr, runtimeErr)
	}

	compileErr, runtimeErr = classifyFailure("python", nil, "", "Traceback", "", "")
	if compileErr != "" || runtimeErr == "" {
		t.Fatalf("unexpected classification for python runtime: compile=%q runtime=%q", compileErr, runtimeErr)
	}
}
