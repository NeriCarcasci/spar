package repo

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

type SyncResult struct {
	Updated bool
	Err     error
}

func PullAsync(repoPath string) <-chan SyncResult {
	ch := make(chan SyncResult, 1)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		result := pull(ctx, repoPath)
		ch <- result
	}()
	return ch
}

func pull(ctx context.Context, repoPath string) SyncResult {
	cmd := exec.CommandContext(ctx, "git", "-C", repoPath, "pull", "--ff-only")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return SyncResult{Err: fmt.Errorf("git pull: %s: %w", string(output), err)}
	}

	alreadyUpToDate := string(output) == "Already up to date.\n" ||
		string(output) == "Already up-to-date.\n"

	return SyncResult{Updated: !alreadyUpToDate}
}
