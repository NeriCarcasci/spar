package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const installModule = "github.com/NeriCarcasci/spar/cmd/spar@latest"

func RunUpdate() {
	fmt.Println(hintStyle.Render("Checking for updates..."))

	goPath, err := exec.LookPath("go")
	if err != nil {
		PrintError("Go toolchain not found. Install Go from https://go.dev/dl/ and try again.")
		os.Exit(1)
	}

	cmd := exec.Command(goPath, "install", installModule)
	cmd.Env = append(os.Environ(), "GOPROXY=direct", "GONOSUMDB=*")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		PrintError("Update failed: %v", err)
		os.Exit(1)
	}

	fmt.Println()
	PrintSuccess("spar updated to latest version")

	version := getInstalledVersion()
	if version != "" {
		fmt.Printf("  %s\n", dimStyle.Render("Installed: "+version))
	}
}

func getInstalledVersion() string {
	sparPath, err := exec.LookPath("spar")
	if err != nil {
		return ""
	}
	cmd := exec.Command("go", "version", "-m", sparPath)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "mod") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				return fields[2]
			}
		}
	}
	return ""
}
