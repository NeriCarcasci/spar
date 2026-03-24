package friends

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type FriendMeta struct {
	Status    string    `json:"status"`
	FetchedAt time.Time `json:"fetched_at"`
}

type SyncMeta struct {
	LastSync time.Time             `json:"last_sync"`
	Results  map[string]FriendMeta `json:"results"`
}

func CacheDir(dataDir string) string {
	return filepath.Join(dataDir, "friends")
}

func LoadCached(dataDir string, username string) (PublicProfile, error) {
	path := filepath.Join(CacheDir(dataDir), username+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return PublicProfile{}, err
	}
	var p PublicProfile
	if err := json.Unmarshal(data, &p); err != nil {
		return PublicProfile{}, fmt.Errorf("parsing cached profile: %w", err)
	}
	return p, nil
}

func SaveCached(dataDir string, username string, profile PublicProfile) error {
	dir := CacheDir(dataDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, username+".json"), data, 0o644)
}

func LoadMeta(dataDir string) (SyncMeta, error) {
	path := filepath.Join(CacheDir(dataDir), "_meta.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return SyncMeta{Results: make(map[string]FriendMeta)}, nil
		}
		return SyncMeta{}, err
	}
	var m SyncMeta
	if err := json.Unmarshal(data, &m); err != nil {
		return SyncMeta{}, err
	}
	if m.Results == nil {
		m.Results = make(map[string]FriendMeta)
	}
	return m, nil
}

func SaveMeta(dataDir string, meta SyncMeta) error {
	dir := CacheDir(dataDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "_meta.json"), data, 0o644)
}
