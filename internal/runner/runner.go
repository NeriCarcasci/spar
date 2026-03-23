package runner

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

type TestResult struct {
	Index    int
	Passed   bool
	Got      string
	Expected string
}

type languageSpec struct {
	BuilderFile  string
	SolutionFile string
}

func Run(challengeDir, language, userSolution string) ([]TestResult, string, string, error) {
	spec, err := specForLanguage(language)
	if err != nil {
		return nil, "", "", err
	}

	tmpDir, err := os.MkdirTemp("", "spar-run-*")
	if err != nil {
		return nil, "", "", err
	}
	defer os.RemoveAll(tmpDir)

	solutionPath := filepath.Join(tmpDir, spec.SolutionFile)
	if err := os.WriteFile(solutionPath, []byte(userSolution), 0o644); err != nil {
		return nil, "", "", err
	}

	builderSrc := filepath.Join(challengeDir, "builder", spec.BuilderFile)
	builderDst := filepath.Join(tmpDir, spec.BuilderFile)
	if err := copyFile(builderSrc, builderDst); err != nil {
		return nil, "", "", err
	}

	testsPath := filepath.Join(challengeDir, "tests.yaml")
	challengePath := filepath.Join(challengeDir, "challenge.yaml")
	stdout, stderr, runErr, compileErr, runtimeErr := execute(tmpDir, solutionPath, testsPath, challengePath, normalizeLanguage(language), spec)
	results := ParseResults(stdout)

	if runErr != nil {
		cErr, rErr := classifyFailure(normalizeLanguage(language), results, stdout, stderr, compileErr, runtimeErr)
		return results, cErr, rErr, nil
	}

	return results, "", "", nil
}

func ParseResults(stdout string) []TestResult {
	passPattern := regexp.MustCompile(`^PASS\s+(\d+)$`)
	failPattern := regexp.MustCompile(`^FAIL\s+(\d+)\s+got=(.+)\s+expected=(.+)$`)
	lines := strings.Split(stdout, "\n")
	results := make([]TestResult, 0)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if matches := passPattern.FindStringSubmatch(trimmed); len(matches) == 2 {
			results = append(results, TestResult{Index: atoi(matches[1]), Passed: true})
			continue
		}
		if matches := failPattern.FindStringSubmatch(trimmed); len(matches) == 4 {
			results = append(results, TestResult{Index: atoi(matches[1]), Passed: false, Got: matches[2], Expected: matches[3]})
		}
	}

	return results
}

func execute(tmpDir, solutionPath, testsPath, challengePath, language string, spec languageSpec) (string, string, error, string, string) {
	switch language {
	case "python":
		stdout, stderr, runErr, compileErr, runtimeErr := runCmd(tmpDir, "python3", spec.BuilderFile, solutionPath, testsPath, challengePath)
		if runErr != nil && isCommandMissing(runErr) {
			return runCmd(tmpDir, "python", spec.BuilderFile, solutionPath, testsPath, challengePath)
		}
		return stdout, stderr, runErr, compileErr, runtimeErr
	case "go":
		return runCmd(tmpDir, "go", "run", spec.BuilderFile, spec.SolutionFile, testsPath, challengePath)
	case "javascript":
		return runCmd(tmpDir, "node", spec.BuilderFile, solutionPath, testsPath, challengePath)
	case "cpp":
		binary := executableName("runner")
		cOut, cErr, cRunErr, _, _ := runCmd(tmpDir, "g++", "-o", binary, spec.BuilderFile)
		if cRunErr != nil {
			return cOut, cErr, cRunErr, strings.TrimSpace(joinNonEmpty(cOut, cErr)), ""
		}
		rOut, rErr, rRunErr, _, _ := runCmd(tmpDir, filepath.Join(tmpDir, binary), solutionPath, testsPath, challengePath)
		return rOut, rErr, rRunErr, "", strings.TrimSpace(joinNonEmpty(rOut, rErr))
	case "rust":
		binary := executableName("runner")
		cOut, cErr, cRunErr, _, _ := runCmd(tmpDir, "rustc", spec.BuilderFile, "-o", binary)
		if cRunErr != nil {
			return cOut, cErr, cRunErr, strings.TrimSpace(joinNonEmpty(cOut, cErr)), ""
		}
		rOut, rErr, rRunErr, _, _ := runCmd(tmpDir, filepath.Join(tmpDir, binary), solutionPath, testsPath, challengePath)
		return rOut, rErr, rRunErr, "", strings.TrimSpace(joinNonEmpty(rOut, rErr))
	default:
		return "", "", fmt.Errorf("unsupported language: %s", language), "", ""
	}
}

func runCmd(dir, name string, args ...string) (string, string, error, string, string) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err, "", ""
}

func isCommandMissing(err error) bool {
	var execErr *exec.Error
	return errors.As(err, &execErr) && errors.Is(execErr, exec.ErrNotFound)
}

func classifyFailure(language string, results []TestResult, stdout, stderr, compileErr, runtimeErr string) (string, string) {
	if compileErr != "" || runtimeErr != "" {
		return compileErr, runtimeErr
	}
	details := strings.TrimSpace(joinNonEmpty(stdout, stderr))
	if details == "" {
		details = "command failed"
	}

	switch language {
	case "python", "javascript":
		return "", details
	case "go":
		if len(results) == 0 {
			return details, ""
		}
		return "", details
	default:
		return "", details
	}
}

func specForLanguage(language string) (languageSpec, error) {
	switch normalizeLanguage(language) {
	case "python":
		return languageSpec{BuilderFile: "builder.py", SolutionFile: "python.py"}, nil
	case "go":
		return languageSpec{BuilderFile: "builder.go", SolutionFile: "solution.go"}, nil
	case "javascript":
		return languageSpec{BuilderFile: "builder.js", SolutionFile: "solution.js"}, nil
	case "cpp":
		return languageSpec{BuilderFile: "builder.cpp", SolutionFile: "solution.cpp"}, nil
	case "rust":
		return languageSpec{BuilderFile: "builder.rs", SolutionFile: "solution.rs"}, nil
	default:
		return languageSpec{}, fmt.Errorf("unsupported language: %s", language)
	}
}

func normalizeLanguage(language string) string {
	switch strings.ToLower(strings.TrimSpace(language)) {
	case "js", "javascript":
		return "javascript"
	case "c++", "cpp":
		return "cpp"
	default:
		return strings.ToLower(strings.TrimSpace(language))
	}
}

func executableName(base string) string {
	if runtime.GOOS == "windows" {
		return base + ".exe"
	}
	return base
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}

func joinNonEmpty(values ...string) string {
	parts := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return strings.Join(parts, "\n")
}

func atoi(value string) int {
	n, _ := strconv.Atoi(value)
	return n
}
