package friends

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spar-cli/spar/internal/challenge"
	"github.com/spar-cli/spar/internal/profile"
	"github.com/spar-cli/spar/internal/rank"
)

func BuildPublicProfile(p *profile.Profile, idx *challenge.Index, version string) PublicProfile {
	solves := map[string]int{"easy": 0, "medium": 0, "hard": 0}
	langs := map[string]int{}
	solved := map[string]bool{}

	for _, s := range p.Solves {
		if !s.Passed {
			continue
		}
		if solved[s.ChallengeID] {
			continue
		}
		solved[s.ChallengeID] = true
		langs[s.Language]++
	}

	if idx != nil {
		for _, e := range idx.Challenges {
			if solved[e.ID] {
				solves[string(e.Difficulty)]++
			}
		}
	}

	trackMedals := map[string]string{}
	for _, name := range rank.TrackNames {
		key := trackKey(name)
		medal, ok := p.TrackMedals[name]
		if !ok {
			medal = 0
		}
		trackMedals[key] = medalString(rank.MedalTier(medal))
	}

	totalChallenges := 0
	if idx != nil {
		totalChallenges = len(idx.Challenges)
	}

	return PublicProfile{
		Username:        p.Username,
		Rank:            p.CurrentTier,
		Division:        p.CurrentDivision,
		TotalSP:         p.TotalSP,
		Streak:          p.Streak,
		Solves:          solves,
		TotalSolved:     len(solved),
		TotalChallenges: totalChallenges,
		Languages:       langs,
		TrackMedals:     trackMedals,
		LastUpdated:     time.Now().UTC(),
		SparVersion:     version,
	}
}

func Publish(repoPath string, remoteName string, pub PublicProfile) error {
	if remoteName == "" {
		remoteName = "origin"
	}

	remoteURL, err := gitRemoteURL(repoPath, remoteName)
	if err != nil {
		return fmt.Errorf("no git remote found. Fork spar on GitHub and add it as a remote")
	}
	if strings.TrimSpace(remoteURL) == "" {
		return fmt.Errorf("no git remote found. Fork spar on GitHub and add it as a remote")
	}

	data, err := json.MarshalIndent(pub, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding profile: %w", err)
	}
	data = append(data, '\n')

	tmpDir, err := os.MkdirTemp("", "spar-publish-*")
	if err != nil {
		return fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	profilePath := filepath.Join(tmpDir, "profile.json")
	if err := os.WriteFile(profilePath, data, 0o644); err != nil {
		return fmt.Errorf("writing temp profile: %w", err)
	}

	exists, err := ProfileBranchExists(repoPath, remoteName)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if exists {
		return publishToExisting(ctx, repoPath, remoteName, tmpDir, profilePath)
	}
	return publishNewBranch(ctx, repoPath, remoteName, tmpDir, profilePath)
}

func ProfileBranchExists(repoPath string, remoteName string) (bool, error) {
	cmd := exec.Command("git", "-C", repoPath, "ls-remote", "--heads", remoteName, "profile")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("checking remote: %w", err)
	}
	return strings.Contains(string(out), "refs/heads/profile"), nil
}

func publishNewBranch(ctx context.Context, repoPath, remoteName, tmpDir, profilePath string) error {
	wtPath := filepath.Join(tmpDir, "worktree")

	if err := gitCmd(ctx, repoPath, "worktree", "add", "--orphan", wtPath, "profile"); err != nil {
		return fmt.Errorf("creating worktree: %w", err)
	}
	defer gitCmd(context.Background(), repoPath, "worktree", "remove", "--force", wtPath)

	entries, _ := os.ReadDir(wtPath)
	for _, e := range entries {
		if e.Name() == ".git" {
			continue
		}
		os.RemoveAll(filepath.Join(wtPath, e.Name()))
	}

	data, _ := os.ReadFile(profilePath)
	if err := os.WriteFile(filepath.Join(wtPath, "profile.json"), data, 0o644); err != nil {
		return fmt.Errorf("writing profile to worktree: %w", err)
	}

	if err := gitCmd(ctx, wtPath, "add", "profile.json"); err != nil {
		return fmt.Errorf("staging: %w", err)
	}
	if err := gitCmd(ctx, wtPath, "commit", "-m", "publish profile"); err != nil {
		return fmt.Errorf("committing: %w", err)
	}
	if err := gitCmd(ctx, repoPath, "push", remoteName, "profile"); err != nil {
		return fmt.Errorf("pushing: %w", err)
	}
	return nil
}

func publishToExisting(ctx context.Context, repoPath, remoteName, tmpDir, profilePath string) error {
	if err := gitCmd(ctx, repoPath, "fetch", remoteName, "profile"); err != nil {
		return fmt.Errorf("fetching profile branch: %w", err)
	}

	wtPath := filepath.Join(tmpDir, "worktree")

	if err := gitCmd(ctx, repoPath, "worktree", "add", wtPath, "profile"); err != nil {
		if err2 := gitCmd(ctx, repoPath, "worktree", "add", "--track", "-b", "profile", wtPath, remoteName+"/profile"); err2 != nil {
			return fmt.Errorf("creating worktree: %w (also tried: %v)", err2, err)
		}
	}
	defer gitCmd(context.Background(), repoPath, "worktree", "remove", "--force", wtPath)

	data, _ := os.ReadFile(profilePath)
	if err := os.WriteFile(filepath.Join(wtPath, "profile.json"), data, 0o644); err != nil {
		return fmt.Errorf("writing profile: %w", err)
	}

	if err := gitCmd(ctx, wtPath, "add", "profile.json"); err != nil {
		return fmt.Errorf("staging: %w", err)
	}

	statusOut, _ := exec.CommandContext(ctx, "git", "-C", wtPath, "status", "--porcelain").Output()
	if len(strings.TrimSpace(string(statusOut))) == 0 {
		return nil
	}

	if err := gitCmd(ctx, wtPath, "commit", "-m", "update profile"); err != nil {
		return fmt.Errorf("committing: %w", err)
	}
	if err := gitCmd(ctx, repoPath, "push", remoteName, "profile"); err != nil {
		return fmt.Errorf("pushing: %w", err)
	}
	return nil
}

func gitRemoteURL(repoPath, remoteName string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "remote", "get-url", remoteName)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func gitCmd(ctx context.Context, dir string, args ...string) error {
	cmd := exec.CommandContext(ctx, "git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func trackKey(name string) string {
	s := strings.ToLower(name)
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, "the_", "")
	return s
}

func medalString(m rank.MedalTier) string {
	switch m {
	case rank.MedalGold:
		return "gold"
	case rank.MedalSilver:
		return "silver"
	case rank.MedalBronze:
		return "bronze"
	default:
		return "none"
	}
}
