package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

var renameMap = map[string]string{
	"solution.py":  "python.py",
	"solution.go":  "go.go",
	"solution.js":  "javascript.js",
	"solution.cpp": "cpp.cpp",
	"solution.rs":  "rust.rs",
}

func main() {
	root := flag.String("root", "challenges", "challenge root directory")
	flag.Parse()

	dirs, err := challengeDirs(*root)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	renamed := 0
	removed := 0
	for _, dir := range dirs {
		r, d, err := normalizeChallenge(dir)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		renamed += r
		removed += d
	}

	fmt.Printf("normalized=%d renamed=%d removed=%d\n", len(dirs), renamed, removed)
}

func challengeDirs(root string) ([]string, error) {
	dirs := []string{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		if entry.Name() == "challenge.yaml" {
			dirs = append(dirs, filepath.Dir(path))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return dirs, nil
}

func normalizeChallenge(challengeDir string) (int, int, error) {
	renamed := 0
	removed := 0
	for _, sub := range []string{"setup", "solutions"} {
		dir := filepath.Join(challengeDir, sub)
		info, err := os.Stat(dir)
		if err != nil || !info.IsDir() {
			continue
		}
		for from, to := range renameMap {
			src := filepath.Join(dir, from)
			dst := filepath.Join(dir, to)
			srcInfo, err := os.Stat(src)
			if err != nil || srcInfo.IsDir() {
				continue
			}
			if _, err := os.Stat(dst); err == nil {
				same, err := sameFile(src, dst)
				if err != nil {
					return renamed, removed, err
				}
				if same {
					if err := os.Remove(src); err != nil {
						return renamed, removed, err
					}
					removed++
				}
				continue
			}
			if err := os.Rename(src, dst); err != nil {
				return renamed, removed, err
			}
			renamed++
		}
	}
	return renamed, removed, nil
}

func sameFile(a, b string) (bool, error) {
	left, err := os.ReadFile(a)
	if err != nil {
		return false, err
	}
	right, err := os.ReadFile(b)
	if err != nil {
		return false, err
	}
	return bytes.Equal(left, right), nil
}
