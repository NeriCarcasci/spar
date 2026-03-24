package friends

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Friend struct {
	URL      string    `yaml:"url"`
	Username string    `yaml:"username"`
	RepoName string    `yaml:"repo_name"`
	AddedAt  time.Time `yaml:"added"`
}

type friendsFile struct {
	Friends []Friend `yaml:"friends"`
}

func LoadFriends(path string) ([]Friend, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading friends: %w", err)
	}
	var f friendsFile
	if err := yaml.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("parsing friends: %w", err)
	}
	return f.Friends, nil
}

func SaveFriends(path string, friends []Friend) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating friends directory: %w", err)
	}
	f := friendsFile{Friends: friends}
	data, err := yaml.Marshal(&f)
	if err != nil {
		return fmt.Errorf("encoding friends: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

func AddFriend(path string, raw string, selfRemoteURL string) (Friend, error) {
	friends, err := LoadFriends(path)
	if err != nil {
		return Friend{}, err
	}

	ghURL, user, repo, err := ParseGitHubURL(raw)
	if err != nil {
		return Friend{}, err
	}

	if selfRemoteURL != "" {
		selfNorm := normalizeRemoteURL(selfRemoteURL)
		if selfNorm == ghURL {
			return Friend{}, fmt.Errorf("cannot add yourself as a friend")
		}
	}

	for _, f := range friends {
		if f.URL == ghURL || strings.EqualFold(f.Username, user) {
			return Friend{}, fmt.Errorf("%s is already in your friends list", user)
		}
	}

	friend := Friend{
		URL:      ghURL,
		Username: user,
		RepoName: repo,
		AddedAt:  time.Now().UTC(),
	}

	friends = append(friends, friend)
	if err := SaveFriends(path, friends); err != nil {
		return Friend{}, err
	}
	return friend, nil
}

func RemoveFriend(path string, identifier string) error {
	friends, err := LoadFriends(path)
	if err != nil {
		return err
	}

	found := -1
	for i, f := range friends {
		if strings.EqualFold(f.Username, identifier) || f.URL == identifier {
			found = i
			break
		}
	}
	if found < 0 {
		return fmt.Errorf("friend %q not found", identifier)
	}

	friends = append(friends[:found], friends[found+1:]...)
	return SaveFriends(path, friends)
}

func FindFriend(friends []Friend, username string) (Friend, bool) {
	for _, f := range friends {
		if strings.EqualFold(f.Username, username) {
			return f, true
		}
	}
	return Friend{}, false
}

func ParseGitHubURL(raw string) (normalized, user, repo string, err error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", "", "", fmt.Errorf("empty URL")
	}

	raw = strings.TrimSuffix(raw, ".git")

	if !strings.Contains(raw, "/") {
		user = raw
		repo = "spar"
		normalized = "https://github.com/" + user + "/" + repo
		return normalized, user, repo, nil
	}

	if !strings.Contains(raw, "://") && !strings.HasPrefix(raw, "github.com") {
		parts := strings.SplitN(raw, "/", 2)
		if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
			user = parts[0]
			repo = parts[1]
			normalized = "https://github.com/" + user + "/" + repo
			return normalized, user, repo, nil
		}
	}

	if strings.HasPrefix(raw, "github.com/") {
		raw = "https://" + raw
	}

	u, parseErr := url.Parse(raw)
	if parseErr != nil || u.Host == "" {
		return "", "", "", fmt.Errorf("invalid URL: %s", raw)
	}

	if !strings.EqualFold(u.Host, "github.com") {
		return "", "", "", fmt.Errorf("invalid URL. Expected a GitHub repository URL")
	}

	pathParts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(pathParts) < 2 || pathParts[0] == "" || pathParts[1] == "" {
		return "", "", "", fmt.Errorf("invalid URL. Expected https://github.com/{user}/{repo}")
	}

	user = pathParts[0]
	repo = pathParts[1]
	normalized = "https://github.com/" + user + "/" + repo
	return normalized, user, repo, nil
}

func normalizeRemoteURL(remote string) string {
	remote = strings.TrimSpace(remote)
	remote = strings.TrimSuffix(remote, ".git")

	if strings.HasPrefix(remote, "git@github.com:") {
		remote = strings.TrimPrefix(remote, "git@github.com:")
		remote = "https://github.com/" + remote
	}

	if strings.HasPrefix(remote, "github.com/") {
		remote = "https://" + remote
	}

	return remote
}
